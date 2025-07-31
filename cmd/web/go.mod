module webserver

go 1.24

toolchain go1.24.0

replace github.com/rorycl/cexfind => ../../

replace github.com/rorycl/cexfind/location => ../../location

require (
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/schema v1.4.1
	github.com/rorycl/cexfind v0.2.8
	github.com/shopspring/decimal v1.4.0
)

require (
	github.com/felixge/httpsnoop v1.0.4 // indirect
	golang.org/x/text v0.27.0 // indirect
)
