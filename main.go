package main

import (
	"bytes"
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
	rootPath  string
	pastePath string
	tmplPath  string
	manTitle  string
	port      string
	formVal   string
	proto     string = "http://"
	manCache  map[string]string
	typeCache map[string]FileType
	maxB      int64
)

// Webpage template
type Template struct {
	Key  string
	Body []byte
	Lang string
}

type FileType uint8

const (
	TypeNil FileType = iota
	TypePlain
	TypePNG
	TypeJPEG
)

// Host a pastebin-like service
func main() {
	flag.StringVar(&rootPath, "r", "./", "Website root directory")
	flag.StringVar(&port, "p", ":8001", "Web server port to host on")
	flag.StringVar(&formVal, "v", "paste", "Form value that appears in 'paste=<-' style form values")
	flag.StringVar(&manTitle, "m", "isepaste", "Title of man page printed on landing page")
	flag.Int64Var(&maxB, "s", 10000000, "Max file size in bytes")
	flag.Parse()

	pastePath = rootPath + "/pastes/"
	tmplPath = rootPath + "/static/"
	manCache = make(map[string]string)
	typeCache = make(map[string]FileType)

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
		manCache[r.Host] = fmt.Sprintf(man, strings.ToLower(manTitle), strings.ToUpper(manTitle), strings.ToLower(manTitle), strings.ToLower(manTitle), formVal, url, formVal, url, url, url, url, formVal, url, url, formVal, "`", url, url, url)
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
	var file string
	var ext string

	if strings.Contains(key, ".") {
		spl := strings.Split(key, ".")
		file = pastePath + spl[0] + ".paste"
		ext = spl[1]
	} else {
		file = pastePath + key + ".paste"
	}

	paste, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(w, "[%s] not found", key)
		return
	}

	if ext == "" {
		// check MIME type
		var mime FileType
		if typeCache[key] == TypeNil {
			b := bytes.NewBuffer(paste)
			mime, err = getFileType(b)
			if err != nil {
				fmt.Fprintf(w, "%s", err)
				return
			}
			typeCache[key] = mime
		} else {
			mime = typeCache[key]
		}

		// redirect plain
		switch mime {
		case TypePNG:
			http.Redirect(w, r, proto+r.Host+"/"+key+".png", 302)
			return
		case TypeJPEG:
			http.Redirect(w, r, proto+r.Host+"/"+key+".jpeg", 302)
			return
		case TypePlain:
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			fmt.Fprintf(w, "%s", paste)
			return
		}
	} else {
		switch ext {
		case "png":
			w.Header().Set("Content-Type", "image/png")
		case "jpeg", "jpg":
			w.Header().Set("Content-Type", "image/jpeg")
		}

		fmt.Fprintf(w, "%s", paste)
		return
	}
}

func getFileType(b *bytes.Buffer) (FileType, error) {
	res, err := isPNG(b)
	if err != nil {
		return TypeNil, err
	}
	if res {
		return TypePNG, nil
	}

	res, err = isJPEG(b)
	if err != nil {
		return TypeNil, err
	}
	if res {
		return TypeJPEG, nil
	}

	return TypePlain, nil
}

func isPNG(buf *bytes.Buffer) (bool, error) {
	return bytes.Compare(buf.Bytes()[0:8], []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) == 0, nil
}
func isJPEG(buf *bytes.Buffer) (bool, error) {
	head := buf.Bytes()[0:2]
	tail := buf.Bytes()[buf.Len()-2 : buf.Len()]
	return (bytes.Compare(head, []byte{0xFF, 0xD8}) == 0) && (bytes.Compare(tail, []byte{0xFF, 0xD9}) == 0), nil
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
