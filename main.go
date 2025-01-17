package main

import (
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Global variables
var (
	rootPath	string
	pastePath	string
	tmplPath	string
	manTitle	string
	port		string
	formVal		string
	proto		string = "http://"
	manCache	map[string]string
	maxB		int64
)

// Webpage template
type Template struct {
	Key  string
	Body []byte
	Lang string
}

// Host a pastebin-like service
func main() {
	flag.StringVar(&rootPath, "r", "./", "Website root directory")
	flag.StringVar(&port, "p", ":8001", "Web server port to host on")
	flag.StringVar(&formVal, "v", "paste", "Form value that appears in 'paste=<-' style form values")
	flag.StringVar(&manTitle, "m", "isepaste", "Title of man page printed on landing page")
	flag.Int64Var(&maxB, "s", 10000000, "Max file size in bytes")
	flag.Parse()

	pastePath	= rootPath + "/pastes/"
	tmplPath	= rootPath + "/static/"
	manCache	= make(map[string]string)

	r := mux.NewRouter()

	// Landing on homepage
	r.HandleFunc("/", handleLand).Methods("GET")

	// Posting a paste
	r.HandleFunc("/", handlePaste).Methods("POST")

	// Reading a paste
	r.HandleFunc("/{pasteId}", handleView).Methods("GET")

	http.Handle("/", r)

	log.Printf("Listening on tcp!*!%s.\n", port[1:])
	log.Fatal(http.ListenAndServe(port, nil))
}

// Landing page handler
func handleLand(w http.ResponseWriter, r *http.Request) {
	if manCache[r.Host] == "" {		
		url := proto + r.Host
		manCache[r.Host] = fmt.Sprintf(man, strings.ToLower(manTitle), strings.ToUpper(manTitle), strings.ToLower(manTitle), strings.ToLower(manTitle), formVal, url, formVal, url, url, url, url, formVal, url, url, formVal, "`", url, url, url )
	}

	fmt.Fprint(w, manCache[r.Host])
}

// Paste path handler — for writing
func handlePaste(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxB)

	paste := r.FormValue(formVal)

	// Generate hash to use as filename/key
	// hash is base64 encoding of the first 72 bits of sha1(paste)
	h := sha1.New()
	h.Write([]byte(paste))
	keyHash := h.Sum(nil)
	key := base64.URLEncoding.EncodeToString(keyHash[:9])

	// Save our paste
	f := pastePath + key + ".paste"
	err := ioutil.WriteFile(f, []byte(paste), 0600)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	u := proto + r.Host + "/" + key
	fmt.Fprintf(w, "%s\n", u)
}

// View path handler — for reading
func handleView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["pasteId"]
	file := pastePath + key + ".paste"

	paste, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(w, "[%s] not found", key)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%s", paste)
	return
}

// Manual for port landing page printing
const man string = `%s(1)                          %s                          %s(1)

NAME
	%s: command line pastebin.

SYNOPSIS
	<command> | curl -F '%s=<-' %s/

DESCRIPTION
	Paste to a listening plaintext paste server.

EXAMPLES
	Paste the file bin/myscript and open the link in firefox(1) from unix:

		~$ cat bin/myscript | curl -F '%s=<-' %s
		%s/aXZI
		~$ firefox %s/aXZI

	Paste the file bin/rc/myscript and plumb the link from Plan 9:

		%% cat bin/rc/myscript | hpost -u %s -p / %s@/fd/0
		%s/aXZI
		%% plumb %s/aXZI

	Paste the file dis/myscript and plumb the link from Inferno:

		; cat dis/myscript | { webgrab -p '%s='^%s{cat /fd/0} -o - %s }
		%s/aXZI
		; plumb %s/aXZI

SOURCE
	https://github.com/henesy/gopaste
`
