package main

import (
	"encoding/json"
	"github.com/gocode/epubviewer/epub"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	logger *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "epubviewer ", log.Lshortfile)

	http.HandleFunc("/", spineHandler)
	http.HandleFunc("/static/", staticFilesHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/toc", tocHandler)
	http.HandleFunc("/epubviewer/", epubViewerHandler)
	http.HandleFunc("/nextpage", nextPageHandler)
	http.HandleFunc("/prevpage", prevPageHandler)
	http.HandleFunc("/filelist", filelistHandler)
}

func uploadHandler(rw http.ResponseWriter, req *http.Request) {

	req.ParseMultipartForm(32 << 20)

	srcFile, header, err := req.FormFile("epubupload")
	if err != nil {
		logger.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	dstFileName := strings.Replace(header.Filename, " ", "_", -1)

	e := epub.New(dstFileName)

	if err := e.Load(srcFile); err != nil {
		logger.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	type nextPage struct {
		Href string
	}

	page := "/epubviewer/" + dstFileName + e.GetToc()
	buf, _ := json.Marshal(nextPage{page})
	rw.Write(buf)
}

func tocHandler(rw http.ResponseWriter, req *http.Request) {
	e := epub.New(req.FormValue("bookname"))

	if err := e.LoadFromCache(); err != nil {
		logger.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	e.WriteToc(rw)
}

func nextPageHandler(rw http.ResponseWriter, req *http.Request) {
	e := epub.New(req.FormValue("bookname"))

	if err := e.LoadFromCache(); err != nil {
		logger.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	page := e.GetNextPage(req.FormValue("href")[1:])

	type nextPage struct {
		Href string
	}

	buf, _ := json.Marshal(nextPage{page})
	rw.Write(buf)
}

func prevPageHandler(rw http.ResponseWriter, req *http.Request) {
	e := epub.New(req.FormValue("bookname"))

	if err := e.LoadFromCache(); err != nil {
		logger.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	c, _ := req.Cookie("bookname")
	e := epub.New(c.Value)

	if err := e.LoadFromCache(); err != nil {
		logger.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	e.WriteFile(rw, req.URL.Path)
}

func filelistHandler(rw http.ResponseWriter, req *http.Request) {
	type file struct {
		Name string
		Href string
	}

	response := make([]file, len(epub.Cache))

	i := 0
	for name, _ := range epub.Cache {
		e := epub.New(name)
		response[i].Name = name
		response[i].Href = "/epubviewer/" + name + e.GetToc()
		i++
	}

	enc := json.NewEncoder(rw)
	enc.Encode(response)
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, "static/index.html")
}

func epubViewerHandler(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, "static/view.html")
}

func staticFilesHandler(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, req.URL.Path[1:])
}

func main() {
	http.ListenAndServe(":8080", nil)
}
