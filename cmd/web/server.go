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
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	cex "github.com/rorycl/cexfind"
)

var (
	// WebMaxHeaderBytes is the largest number of header bytes accepted by
	// the webserver
	WebMaxHeaderBytes int = 1 << 17 // ~125k

	// ServerAddress is the default Server network address
	ServerAddress string = "127.0.0.1"

	// ServerPort is the default Server network port
	ServerPort string = "8000"

	// BaseURL is the base url for redirects, etc.
	BaseURL string = ""
)

// searcher is an indirect of cex.Search to allow testing
var searcher func(queries []string, strict bool) ([]cex.Box, error) = cex.Search

// listenAndServe is an indirect of http/net.Server.ListenAndServe
var listenAndServe = (*http.Server).ListenAndServe

// development flags and static and template directory locations
var (
	// production is default; set inDevelopment to true with build tag
	inDevelopment bool   = false
	staticDirDev  string = "web/static"
	tplDirDev     string = "web/templates"
	staticDir     string = "static"
	tplDir        string = "templates"
	DirFS         *fileSystem
)

// setupFS setup the filesystem for templates or static files, depending on
// development (filesystem) or not (embedded)
func setupFS() error {
	var err error
	if inDevelopment {
		DirFS, err = NewFileSystem(inDevelopment, tplDirDev, staticDirDev)
	} else {
		DirFS, err = NewFileSystem(inDevelopment, tplDir, staticDir)
	}
	return err
}

// Serve runs the web server on the specified address and port
func Serve(addr, port string) {

	if addr == "" {
		addr = ServerAddress
	} else {
		ServerAddress = addr
	}

	if port == "" {
		port = ServerPort
	} else {
		ServerPort = port
	}

	// setup the filesystem
	if err := setupFS(); err != nil {
		log.Fatal(err)
	}

	// endpoint routing; gorilla mux is used because "/" in http.NewServeMux
	// is a catch-all pattern
	r := mux.NewRouter()

	// attach static dynamic file system to the http.FileServer
	// https://pkg.go.dev/github.com/gorilla/mux#section-readme :Static Files
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.FS(DirFS.StaticFS))),
	)

	// routes
	r.HandleFunc("/results", Results)
	r.HandleFunc("/health", Health)
	r.HandleFunc("/favicon.ico", Favicon)
	r.HandleFunc("/", Home)

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
		Addr:    addr + ":" + port,
		Handler: r,
		// timeouts and limits
		MaxHeaderBytes:    WebMaxHeaderBytes,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
	log.Printf("serving on %s:%s", addr, port)

	err := listenAndServe(server)
	if err != nil {
		log.Printf("fatal server error: %v", err)
	}
}

// Results shows the results of a "search" form submission in an htmx partial
func Results(w http.ResponseWriter, r *http.Request) {

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
	type PostResults struct {
		Strict  bool   `schema:"strict"`
		Queries string `schema:"queries"`
	}
	var postResults PostResults
	var decoder = schema.NewDecoder() // best as package decoder
	err = decoder.Decode(&postResults, urlVals)
	if err != nil || len(postResults.Queries) == 0 {
		log.Printf("cex POST : %+v %v", postResults, err)
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, "no query found")
		return
	}

	// split the comma delimited query into queries
	queries := []string{}
	qs := strings.Split(postResults.Queries, ",")
	for _, q := range qs {
		queries = append(queries, strings.TrimSpace(q))
	}

	base := fmt.Sprintf("strict=%s", func() string {
		if postResults.Strict {
			return "true"
		}
		return "false"
	}())
	for _, q := range queries {
		base += fmt.Sprintf("&q=%s", url.PathEscape(q))
	}
	// push the postResults terms to the url
	w.Header().Set("HX-Push-Url", BaseURL+"/?"+base)

	// search; note that searcher is an indirect to search/cex.Search
	type SearchResults struct {
		Results []cex.Box
		Err     error
	}
	sr := SearchResults{}
	sr.Results, sr.Err = searcher(queries, postResults.Strict)

	t := template.Must(template.ParseFS(DirFS.TplFS, "partial-results.html"))
	err = t.Execute(w, sr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "template writing problem : %s", err.Error())
	}
}

// Home is the home page
func Home(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFS(DirFS.TplFS, "home.html")
	if err != nil {
		log.Fatal(err)
	}

	type Search struct {
		Strict  bool     `schema:"strict"`
		Queries []string `schema:"q"`
	}
	var search Search
	var decoder = schema.NewDecoder() // best as package decoder
	err = decoder.Decode(&search, r.URL.Query())

	if inDevelopment {
		log.Printf("cex url GET : %+v %+v (%d items) err %v", r.URL.Query(), search, len(search.Queries), err)
	}

	data := struct {
		Title   string
		Address string
		Port    string
		Search  Search
	}{
		"search cex",
		ServerAddress,
		ServerPort,
		search,
	}
	err = t.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "template writing problem : %s", err.Error())
	}
}

// HealthCheck shows if the service is up
func Health(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := map[string]string{"status": "up"}
	if err := enc.Encode(resp); err != nil {
		log.Print("health error: unable to encode response")
	}
}

// Favicon serves up the favicon
func Favicon(w http.ResponseWriter, r *http.Request) {
	log.Println("got favicon")
	http.ServeFileFS(w, r, DirFS.StaticFS, "/favicon.svg")
}
