# cexfind

v0.2.4 : 21 June 2024 : add buy/exchange and store info; htmx 2.0

A Go module with console, cli and web app clients for rapid and
effective searches for second hand equipment for sale at Cex/Webuy using
the unofficial `webuy.io` json search endpoint.

Note that these programs only work for queries made in the UK (or via a
proxy terminating in the UK). This is intended to be a fun project and
is not intended for commercial use.

The web client is presently deployed on Google Cloud Run. [Try it
out!](https://cexfind-min-poyflf5akq-nw.a.run.app/).

## Usage

Simply download the binaries for your machine's architecture from
[releases](https://github.com/rorycl/cexfind/releases). Alternatively,
build for your local machine using `make build-all` if you have go (>=
1.22) installed. The resulting binaries can be found in `bin`.

## Clients

Three clients are provided for the very simple `cexfind` golang module:

**console**

A [bubbletea](https://github.com/charmbracelet/bubbletea) console app.

<img width="1000" src="cmd/console/console.gif" />

Have a look at the app [README](cmd/console/README.md) for more info
about the architecture of this client.

**cli**

A simple cli client.

Run `./bin/cli -h` or the windows alternative to see the switch options.

<img width="1000" src="cmd/cli/cli.gif" />

**web server**

A simple htmx webserver client.

Run `./bin/webserver` or the windows alternative to run the server
locally on the default local ip address of `127.0.0.1` and port `8000`.
Use the command line switches to change these options. (Use `-h` to see
the switches.)

<img width="1000" src="cmd/web/web.gif" />

## Licence

This project is licensed under the [MIT Licence](LICENCE).
