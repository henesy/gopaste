# GoPaste

__Simple command line pastebin-esque site for personal use.__

Using Gorilla Toolkit [Mux](http://www.gorillatoolkit.org/pkg/mux) for routing.

## Build

 - `go build`

## Usage

`cat file-to-upload.txt | curl -F 'paste=<-' http://your-site`

I like to add ``` | xargs firefox``` to the end to open it in firefox

A plaintext response is served.

## Thanks

Thanks for http://sprunge.us for the idea which I shamelessly copied. Thanks to Gorilla Toolkit for 
the awesome Golang http extensions.

