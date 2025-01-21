package main

// https://ieftimov.com/posts/testing-in-go-testing-http-servers/
// https://bignerdranch.com/blog/using-the-httptest-package-in-golang/

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	cex "github.com/rorycl/cexfind"
	"github.com/shopspring/decimal"
)

// TestSetupFS sets up the FS
func TestSetupFS(t *testing.T) {
	staticDirDev = "static"
	tplDirDev = "templates"
	if got, want := setupFS(), error(nil); got != want {
		t.Errorf("testsetupfs error got %v != want %v", got, want)
	}
}

// TestServe
func TestServe(t *testing.T) {
	listenAndServe = func(*http.Server) error {
		return nil
	}
	staticDirDev = "static"
	tplDirDev = "templates"
	Serve("127.0.0.1", "8123")

}

// Test Home page returns a 200
func TestHome(t *testing.T) {

	// home uses templates fs
	DirFS = &fileSystem{}
	DirFS.TplFS = os.DirFS("templates")

	r := httptest.NewRequest(http.MethodGet, "http://example.com/home", nil)
	w := httptest.NewRecorder()

	Home(w, r)

	res := w.Result()
	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if want, got := 200, res.StatusCode; want != got {
		t.Errorf("expected status %d, got %d", want, got)
	}
}

// Test Health page returns a 200
func TestHealth(t *testing.T) {

	r := httptest.NewRequest(http.MethodGet, "http://example.com/health", nil)
	w := httptest.NewRecorder()

	Health(w, r)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if want, got := 200, res.StatusCode; want != got {
		t.Errorf("expected status %d, got %d", want, got)
	}
	responseBody := string(data)
	if want, got := strings.TrimSpace(`{"status":"up"}`), strings.TrimSpace(responseBody); want != got {
		t.Errorf("expected status %s, got %s", want, got)
	}
}

// Favicon page returns a 200
func TestFavicon(t *testing.T) {

	// favicon uses the static fs
	DirFS = &fileSystem{}
	DirFS.StaticFS = os.DirFS("static")

	r := httptest.NewRequest(http.MethodGet, "http://example.com/favicon.ico", nil)
	w := httptest.NewRecorder()

	Favicon(w, r)

	res := w.Result()
	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if want, got := 200, res.StatusCode; want != got {
		t.Errorf("expected status %d, got %d", want, got)
	}
}

// TestResults tests a POST to Results; note that cex.Search is
// swapped out
func TestResults(t *testing.T) {

	// results uses the templates endpoint
	DirFS = &fileSystem{}
	DirFS.TplFS = os.DirFS("templates")

	// override package global searcher which indirects Search
	searcher = func(queries []string, strict bool, postcode string) ([]cex.Box, error) {
		return []cex.Box{
			cex.Box{Model: "2a", Name: "2a name", ID: "id3", Price: decimal.NewFromInt(3)},
			cex.Box{Model: "1a", Name: "1a name", ID: "id1", Price: decimal.NewFromInt(1)},
			cex.Box{Model: "1b", Name: "1b name", ID: "id2", Price: decimal.NewFromInt(2)},
		}, nil
	}

	tt := []struct {
		name       string
		method     string
		input      string
		statusCode int
	}{
		{
			name:       "succeed post",
			method:     http.MethodPost,
			input:      "query=abc&query=def&strict=false",
			statusCode: http.StatusOK,
		},
		{
			name:       "fail post query too short",
			method:     http.MethodPost,
			input:      "query=ab&strict=false",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "fail post query too short 2",
			method:     http.MethodPost,
			input:      "query=abc&query=de&strict=false",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "fail get",
			method:     http.MethodGet,
			input:      "query=abc&query=def&strict=false",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "fail no POST body",
			method:     http.MethodPost,
			input:      "",
			statusCode: http.StatusNoContent,
		},
	}

	for _, tc := range tt {
		t.Logf("%+v\n", tc)
		t.Run(tc.name, func(t *testing.T) {

			r := httptest.NewRequest(tc.method, "http://example.com/request", strings.NewReader(tc.input))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			Results(w, r)

			res := w.Result()
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			if tc.statusCode != res.StatusCode {
				t.Errorf("expected status %d, got %d", tc.statusCode, res.StatusCode)
			}

		})
	}
}
