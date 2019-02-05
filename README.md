# GoPaste

Simple command line pastebin-esque site for personal use.

## Dependencies

Vendoring Gorilla [Mux](http://www.gorillatoolkit.org/pkg/mux) for routing.

## Build

 - `go build -mod=vendor`

## Usage

Upload myfile.txt and receive the URL back.

	cat myfile.txt | curl -F 'paste=<-' http://your-site

You can add ` | xargs firefox` to the end to open it in firefox, etc.

A plaintext response is served.

The landing page provides a man(1)-style manual page for reference by users.

## Thanks

Thanks for http://sprunge.us for the idea which I shamelessly copied.

Thanks to Gorilla Toolkit for the awesome Golang http extensions.

