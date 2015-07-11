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

func uploadHandler(rw http.ResponseWriter, req *http.Request) {

	req.ParseMultipartForm(32 << 20)

	src, _, err := req.FormFile("epubupload")
	if err != nil {
		logger.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	dst, err := ioutil.TempFile(".files", "file")
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
}

func tocHandler(rw http.ResponseWriter, req *http.Request) {
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
		indexHandler(rw, req)
		return
	}

	e.WriteFile(rw, req.URL.Path)
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, "static/index.html")
}

func staticFilesHandler(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, req.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", spineHandler)
	http.HandleFunc("/static/", staticFilesHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/toc", tocHandler)

	defer func() {
		err := recover()
		log.Println(err)
	}()

	http.ListenAndServe(":9090", nil)
}
