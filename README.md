# cexfind

A Go module with console, cli and web app clients for rapid and
effective searches for equipment on Cex/Webuy using the unofficial
`webuy.io` json search endpoint. Note that these programs only work for
queries made in the UK (or via a proxy terminating in the UK.)

## Usage

Simply download the binaries for your machine's architecture from [the
project releases page](https://github.com/rorycl/cexfind/releases).
Alternatively, build for your local machine using `make build-all` if
you have go (>= 1.22) installed. The resulting binaries can be found in
`bin`.

## Clients

Three clients are provided for the very simple `cexfind` golang module:

**console**

A [bubbletea](https://github.com/charmbracelet/bubbletea) console app.

<img width="1000" src="cmd/console/console.gif" />

Have a look at the [README](cmd/console/README.md) for the console app
for more info about the architecture of this client.

**cli**

A simple cli client.

Run `./bin/cli -h` or the windows alternative to see the switch options.

<img width="1000" src="cmd/web/cli.gif" />

**web server**

A simple htmx webserver client.

Run `./bin/webserver` or the windows alternative to run the server
locally on the default local ip address of `127.0.0.1` and port `8000`.
Use the command line switches to change these options. (Use `-h` to see
the switches.)

<img width="1000" src="cmd/web/web.gif" />


## Licence

This project is licensed under the [MIT Licence](LICENCE).
