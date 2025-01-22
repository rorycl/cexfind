package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestWebMainFlags(t *testing.T) {

	var exit int
	Exit = func(code int) {
		exit = code
	}

	tests := []struct {
		args     []string
		exitCode int
		address  string
		port     string
	}{
		{
			args:     []string{"prog"},
			exitCode: 0,
			address:  "127.0.0.1",
			port:     "8000",
		},
		{
			args:     []string{"prog", "-port", "8002"},
			exitCode: 0,
			address:  "127.0.0.1",
			port:     "8002",
		},
		{
			args:     []string{"prog", "-port", "abc"},
			exitCode: 1,
		},
		{
			args:     []string{"prog", "-address", "127.0.0.2"},
			exitCode: 0,
			address:  "127.0.0.2",
			port:     "8000",
		},
		{
			args:     []string{"prog", "-address", "a.b.c.d"},
			exitCode: 1,
		},
		{
			args:     []string{"prog", "-address", "127.0.0.3", "-port", "8001"},
			exitCode: 0,
			address:  "127.0.0.3",
			port:     "8001",
		},
	}

	for i, tt := range tests {

		// reset the flag environment
		exit = 0
		flag.CommandLine = flag.NewFlagSet(fmt.Sprintf("%d", i), flag.ContinueOnError)

		os.Args = tt.args

		a, p := flagGet()
		t.Logf("subtest %d, args %v", i, tt.args)
		if got, want := exit, tt.exitCode; got != want {
			t.Errorf("got exit code %d expected %d", got, want)
		}
		if exit == tt.exitCode && exit == 0 {
			got, want := a+":"+p, tt.address+":"+tt.port
			if got != want {
				t.Errorf("address/port got %s want %s", got, want)
			}
		}
	}
}

func TestWebMain(t *testing.T) {

	flagGetter = func() (string, string) {
		return "127.0.0.1", "8000"
	}
	if a, b := flagGetter(); a+":"+b != "127.0.0.1:8000" {
		t.Errorf("expected indirected flagGetter == 127.0.0.1:8000")
	}
	serveFunc = func(s *server, address, port string) { log.Print("got here") }
	main()
}
