module webserver

go 1.22.5

replace github.com/rorycl/cexfind => ../../

require (
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/schema v1.4.1
	github.com/rorycl/cexfind v0.2.6
	github.com/shopspring/decimal v1.4.0
)

require (
	github.com/felixge/httpsnoop v1.0.4 // indirect
	golang.org/x/text v0.21.0 // indirect
)
