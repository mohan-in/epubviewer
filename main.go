package main

import (
	"encoding/json"
	"github.com/gocode/epubviewer/epub"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	e      *epub.Ebook
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

	e = epub.New(dst.Name())

	if err := e.Load(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	type page struct {
		Href string
	}

	p := e.GetToc()
	buf, _ := json.Marshal(page{p})
	rw.Write(buf)
}

func tocHandler(rw http.ResponseWriter, req *http.Request) {
	if err := e.Load(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	//e.WriteToc(rw)
}

func nextPageHandler(rw http.ResponseWriter, req *http.Request) {
	page := e.GetNextPage(req.FormValue("href")[1:])

	type nextPage struct {
		Href string
	}

	buf, _ := json.Marshal(nextPage{page})
	rw.Write(buf)
}

func prevPageHandler(rw http.ResponseWriter, req *http.Request) {
	page := e.GetPrevPage(req.FormValue("href")[1:])

	type prevPage struct {
		Href string
	}

	buf, _ := json.Marshal(prevPage{page})
	rw.Write(buf)
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
	http.HandleFunc("/nextpage", nextPageHandler)
	http.HandleFunc("/prevpage", prevPageHandler)

	defer func() {
		err := recover()
		log.Println(err)
	}()

	http.ListenAndServe(":9090", nil)
}
