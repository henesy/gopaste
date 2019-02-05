package main

import (
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"log"
	"./mux"				// https://github.com/gorilla/mux
)


// Global variables
var (
	ROOT_PATH string
	PASTE_PATH string
	TMPLT_PATH string
	MAN_TITLE string

	SITE_URL string
	LISTEN_PORT string

	FORM_VALUE string
)


// Webpage template
type Template struct {
	Key string
	Body []byte
	Lang string
}


/*-------------------------------------
	Main
-------------------------------------*/

func main() {
	flag.StringVar(&ROOT_PATH, "r", "./", "Website root directory")
	flag.StringVar(&SITE_URL, "s", "http://paste.iseage.org", "Website base url for paste output")
	flag.StringVar(&LISTEN_PORT, "p", ":8001", "Web server port to host on")
	flag.StringVar(&FORM_VALUE, "v", "paste", "Form value that appears in 'paste=<-' style form values")
	flag.StringVar(&MAN_TITLE, "m", "isepaste", "Title of man page printed on landing page")
	flag.Parse()

	PASTE_PATH = ROOT_PATH + "/pastes/"
	TMPLT_PATH = ROOT_PATH + "/static/"

	r := mux.NewRouter()

	// Landing on homepage
	r.HandleFunc("/", handleLand).Methods("GET")

	// Posting a paste
	r.HandleFunc("/", handlePaste).Methods("POST")
	
	// Reading a paste
	r.HandleFunc("/{pasteId}", handleView).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(LISTEN_PORT, nil))

	fmt.Println("Good bye! ☺")
}


/*-------------------------------------
	Landing Handler
-------------------------------------*/

func handleLand(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, man, strings.ToLower(MAN_TITLE), strings.ToUpper(MAN_TITLE), strings.ToLower(MAN_TITLE), FORM_VALUE, SITE_URL, LISTEN_PORT, FORM_VALUE, SITE_URL, SITE_URL, SITE_URL)
}


/*-------------------------------------
	Paste Handler
-------------------------------------*/

func handlePaste(w http.ResponseWriter, r *http.Request) {
	paste := r.FormValue(FORM_VALUE)

	// Generate hash to use as filename/key
	// hash is base64 encoding of the first 72 bits of sha1(paste)
	h := sha1.New()
	h.Write([]byte(paste))
	keyHash := h.Sum(nil)
	key := base64.URLEncoding.EncodeToString(keyHash[:9])

	// Save our paste
	f := PASTE_PATH + key + ".paste"
	err := ioutil.WriteFile(f, []byte(paste), 0600)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	u := SITE_URL + "/" + key
	fmt.Fprintf(w, "%s\n", u)
}


/*-------------------------------------
	View Handler	
-------------------------------------*/

func handleView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["pasteId"]
	file := PASTE_PATH + key + ".paste"

	paste, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(w, "[%s] not found", key)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s", paste)
	return
}

// Manual for port landing page printing
var man string = `%s(1)                          %s                          %s(1)

NAME
    isepaste: command line pastebin.

SYNOPSIS
    <command> | curl -F '%s=<-' %s%s/

DESCRIPTION
	Paste to a listening plaintext bin server.

EXAMPLES
    ~$ cat bin/myscript | curl -F '%s=<-' %s
       %s/aXZI
    ~$ firefox %s/aXZI

SEE ALSO
	https://github.com/ISEAGE-ISU/gopaste
`
