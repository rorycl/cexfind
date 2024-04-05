# cexfind

A Go cli and web app for rapid and effective searches for equipment on
Cex/Webuy using the unofficial `webuy.io` json search endpoint.

## Usage

Simply download the binaries for your machine's architecture from [the
project releases page](https://github.com/rorycl/cexfind/releases).
Alternatively, build for your local machine using `make build-all` if
you have go (>= 1.22) installed. The resulting binaries can be found in
`bin`.

**web server**

Run `./bin/webserver` or the windows alternative to run the server
locally on the default local ip address of `127.0.0.1` and port `8000`.
Use the command line switches to change these options. (Use `-h` to see
the switches.)

![](media/web.gif)

**cli**

Run `./bin/cli -h` or the windows alternative to see the switch options.

![](media/cli.gif)

## Licence

This project is licensed under the [MIT Licence](LICENCE).
