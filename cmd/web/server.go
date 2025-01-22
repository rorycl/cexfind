package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	"github.com/rorycl/cexfind"
	cex "github.com/rorycl/cexfind"
	"github.com/rorycl/cexfind/cmd"
)

// listenAndServe is an indirect of http/net.Server.ListenAndServe
var listenAndServe = (*http.Server).ListenAndServe

// production is default; set inDevelopment to true with build tag
var inDevelopment bool = false

// the server struct holds development flags and static and template
// directory locations
type server struct {
	WebMaxHeaderBytes int
	ServerAddress     string
	ServerPort        string
	BaseURL           string

	// cex is the main cexfind plug point
	cex *cexfind.CexFind

	// searcher is an indirect of cex.Search to allow testing
	searcher func(cex *cexfind.CexFind, queries []string, strict bool, postcode string) ([]cex.Box, error)

	staticDirDev string
	tplDirDev    string
	staticDir    string
	tplDir       string
	DirFS        *fileSystem

	// serveFunc is an indirect for the main server functionality,
	// provided for testing
	serveFunc func()
}

func newServer() *server {
	s := server{
		// paths
		staticDirDev: "static",
		tplDirDev:    "templates",
		staticDir:    "static",
		tplDir:       "templates",

		// initialise the search apparatus
		cex: cexfind.NewCexFind(),

		// searcher is an indirect of cex.Search to allow testing
		searcher: (*cex.CexFind).Search,

		// WebMaxHeaderBytes is the largest number of header bytes accepted by
		// the webserver
		WebMaxHeaderBytes: 1 << 17, // ~125k

		// ServerAddress is the default Server network address
		ServerAddress: "127.0.0.1",

		// ServerPort is the default Server network port
		ServerPort: "8000",

		// BaseURL is the base url for redirects, etc.
		BaseURL: "",
	}
	// serveFunc is an indirect for testing
	s.serveFunc = s.serve
	return &s
}

// setupFS setup the filesystem for templates or static files, depending on
// development (filesystem) or not (embedded)
func (s *server) setupFS() error {
	var err error
	if inDevelopment {
		s.DirFS, err = NewFileSystem(inDevelopment, s.tplDirDev, s.staticDirDev)
	} else {
		s.DirFS, err = NewFileSystem(inDevelopment, s.tplDir, s.staticDir)
	}
	return err
}

// Serve runs the web server on the specified address and port
func (s *server) Serve(addr, port string) {
	if addr != "" {
		s.ServerAddress = addr
	}
	if port != "" {
		s.ServerPort = port
	}
	// setup the filesystem
	if err := s.setupFS(); err != nil {
		log.Fatal(err)
	}
	s.serveFunc()
}

func (s *server) serve() {
	// endpoint routing; gorilla mux is used because "/" in http.NewServeMux
	// is a catch-all pattern
	r := mux.NewRouter()

	// attach static dynamic file system to the http.FileServer
	// https://pkg.go.dev/github.com/gorilla/mux#section-readme :Static Files
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.FS(s.DirFS.StaticFS))),
	)

	// routes
	r.HandleFunc("/results", s.Results)
	r.HandleFunc("/health", s.Health)
	r.HandleFunc("/favicon.ico", s.Favicon)
	r.HandleFunc("/", s.Home)

	// logging converts gorilla's handlers.CombinedLoggingHandler to a
	// func(http.Handler) http.Handler to satisfy type MiddlewareFunc
	logging := func(handler http.Handler) http.Handler {
		return handlers.CombinedLoggingHandler(os.Stdout, handler)
	}

	// recovery converts gorilla's handlers.RecoveryHandler to a
	// func(http.Handler) http.Handler to satisfy type MiddlewareFunc
	recovery := func(handler http.Handler) http.Handler {
		return handlers.RecoveryHandler()(handler)
	}

	// compression handler
	compressor := func(handler http.Handler) http.Handler {
		return handlers.CompressHandler(handler)
	}

	// attach middleware
	// r.Use(bodyLimitMiddleware)
	r.Use(logging)
	r.Use(compressor)
	r.Use(recovery)

	// configure server options
	server := &http.Server{
		Addr:    s.ServerAddress + ":" + s.ServerPort,
		Handler: r,
		// timeouts and limits
		MaxHeaderBytes:    s.WebMaxHeaderBytes,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
	log.Printf("serving on %s:%s", s.ServerAddress, s.ServerPort)

	err := listenAndServe(server)
	if err != nil {
		log.Printf("fatal server error: %v", err)
	}
}

// Results shows the results of a "search" form submission in an htmx partial
func (s *server) Results(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("endpoint only accepts POST requests, got", r.Method)
		return
	}

	// read body
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("results endpoint body reading error", err)
		return
	}
	if inDevelopment {
		log.Println("body content:", string(body))
	}

	// extract query from POSTed htmx form
	urlVals, err := url.ParseQuery(string(body))
	if err != nil {
		log.Printf("url parsequery error: %v", err)
		return
	}
	var postResults QueriesType
	var decoder = schema.NewDecoder() // best as package decoder
	err = decoder.Decode(&postResults, urlVals)
	if err != nil || len(postResults.Query) == 0 {
		log.Printf("cex POST : %+v %v", postResults, err)
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, "no query found")
		return
	}

	// split the comma delimited query into queries
	queries, err := cmd.QueryInputChecker(postResults.Query...)
	if err != nil {
		log.Printf("cex queries error: %v %v", postResults.Query, err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "query error: %v", err)
		return
	}

	base := fmt.Sprintf("strict=%s", func() string {
		if postResults.Strict {
			return "true"
		}
		return "false"
	}())
	if postResults.Postcode != "" {
		base += fmt.Sprintf("&postcode=%s", url.PathEscape(postResults.Postcode))
	}
	for _, q := range queries {
		base += fmt.Sprintf("&query=%s", url.PathEscape(q))
	}
	// push the postResults terms to the url
	w.Header().Set("HX-Push-Url", s.BaseURL+"/?"+base)

	// search; note that searcher is an indirect to search/cex.Search
	type SearchResults struct {
		Results []cex.Box
		Err     error
	}
	sr := SearchResults{}
	sr.Results, sr.Err = s.searcher(s.cex, queries, postResults.Strict, postResults.Postcode)

	t := template.Must(template.ParseFS(s.DirFS.TplFS, "partial-results.html"))
	err = t.Execute(w, sr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "template writing problem : %s", err.Error())
	}
}

type QueriesType struct {
	Postcode string   `schema:"postcode"`
	Strict   bool     `schema:"strict"`
	Query    []string `schema:"query"`
}

// String provides a string representation of QueriesType.Query,
// suitable for use in a template
func (q QueriesType) String() string {
	output := ""
	for i, query := range q.Query {
		if i > 0 {
			output += cmd.QuerySplitChar + " "
		}
		output += query
	}
	return output
}

// Home is the home page
func (s *server) Home(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFS(s.DirFS.TplFS, "home.html")
	if err != nil {
		log.Fatal(err)
	}

	var search QueriesType
	var decoder = schema.NewDecoder() // best as package decoder
	err = decoder.Decode(&search, r.URL.Query())

	if inDevelopment {
		log.Printf("cex url GET : %+v %+v (%d items) err %v", r.URL.Query(), search, len(search.Query), err)
	}

	data := struct {
		Title   string
		Address string
		Port    string
		Search  QueriesType
	}{
		"search cex",
		s.ServerAddress,
		s.ServerPort,
		search,
	}
	err = t.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "template writing problem : %s", err.Error())
	}
}

// HealthCheck shows if the service is up
func (s *server) Health(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := map[string]string{"status": "up"}
	if err := enc.Encode(resp); err != nil {
		log.Print("health error: unable to encode response")
	}
}

// Favicon serves up the favicon
func (s *server) Favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, s.DirFS.StaticFS, "/favicon.svg")
}
