package main

import (
	"github.com/gocode/epubviewer/epub"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	e      *epub.Epub
	logger *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "epubviewer ", log.Lshortfile)
}

func uploadHandler(rw http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)

	src, _, err := r.FormFile("epubupload")
	if err != nil {
		logger.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	dst, err := ioutil.TempFile("files", "file")
	if err != nil {
		logger.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	io.Copy(dst, src)

	e, err = epub.New(dst.Name())
	if err != nil {
		logger.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := e.WriteToc(rw); err != nil {
		if err := e.WriteSpine(rw); err != nil {
			logger.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func spineHandler(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		http.ServeFile(rw, req, "index.html")
	} else {
		e.WriteFile(rw, req.URL.Path)
	}
}

func indexHandler(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, "index.html")
}

func main() {
	http.HandleFunc("/", spineHandler)
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/toc", uploadHandler)

	http.ListenAndServe(":9090", nil)

}
