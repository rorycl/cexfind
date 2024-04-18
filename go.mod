module github.com/rorycl/cexfind

go 1.22

replace github.com/rorycl/cexfind/search => ./search

replace github.com/rorycl/cexfind/web => ./web

require (
	github.com/google/go-cmp v0.6.0
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/schema v1.3.0
	golang.org/x/text v0.14.0
)

require github.com/felixge/httpsnoop v1.0.3 // indirect
