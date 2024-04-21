// A webserver client to github.com/rorycl/cexfind
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
)

var Server func(address, port string) = Serve

var usage = `
run a webserver to search Cex/Webuy for second hand equipment

eg <programme> [-a 127.0.0.1] [-p 8001]
`

// indirect Exit for testing
var Exit func(code int) = os.Exit

// flagGetter indirects flagGet for testing
var flagGetter func() (string, string) = flagGet

// flagGet checks the flags
func flagGet() (address, port string) {

	flag.StringVar(&address, "address", "127.0.0.1", "server network address")
	flag.StringVar(&port, "port", "8000", "server network port")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(flag.CommandLine.Output(), usage)
	}

	flag.Parse()
	if address == "" || port == "" {
		flag.Usage()
		Exit(1)
	}

	// check address and port are valid
	a := net.ParseIP(address)
	if a == nil {
		fmt.Printf("address %s invalid\n", address)
		Exit(1)
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		fmt.Printf("port %s invalid\n", port)
		Exit(1)
	}

	return
}

func main() {
	address, port := flagGetter()
	Server(address, port)
}
