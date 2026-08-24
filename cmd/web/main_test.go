package main

import (
	"flag"
	"fmt"
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
		proxy    string
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
		{
			args:     []string{"prog", "-address", "127.0.0.3", "-port", "8001", "-proxy", "socks5://127.0.0.1:8001"},
			exitCode: 0,
			address:  "127.0.0.3",
			port:     "8001",
			proxy:    "socks5://127.0.0.1:8001",
		},
	}

	for i, tt := range tests {

		// reset the flag environment
		exit = 0
		flag.CommandLine = flag.NewFlagSet(fmt.Sprintf("%d", i), flag.ContinueOnError)

		os.Args = tt.args

		a, p, proxy := flagGet()
		t.Logf("subtest %d, args %v", i, tt.args)
		if got, want := exit, tt.exitCode; got != want {
			t.Errorf("got exit code %d expected %d", got, want)
		}
		if exit == tt.exitCode && exit == 0 {
			got, want := a+":"+p+":"+proxy, tt.address+":"+tt.port+":"+tt.proxy
			if got != want {
				t.Errorf("address/port/proxy got %s want %s", got, want)
			}
		}
	}
}
