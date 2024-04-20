module webserver

go 1.22

replace github.com/rorycl/cexfind => ../../

require (
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/schema v1.3.0
	github.com/rorycl/cexfind v0.0.0-00010101000000-000000000000
)

require (
	github.com/felixge/httpsnoop v1.0.3 // indirect
	golang.org/x/text v0.14.0 // indirect
)
