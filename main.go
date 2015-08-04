package main

import (
	"encoding/json"
	"github.com/gocode/epubviewer/epub"
	"io"
	"io/ioutil"
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

	e := epub.New(dst.Name())

	if err := e.Load(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	c := http.Cookie{Name: "BookName", Value: strings.Replace(dst.Name(), "\\", "*", -1)}
	http.SetCookie(rw, &c)

	type nextPage struct {
		Href string
	}

	page := "/epubviewer/" + strings.Replace(dst.Name(), "\\", "*", -1) + e.GetToc()
	buf, _ := json.Marshal(nextPage{page})
	rw.Write(buf)
}

func tocHandler(rw http.ResponseWriter, req *http.Request) {
	c, _ := req.Cookie("BookName")
	e := epub.New(strings.Replace(c.Value, "*", "\\", -1))
	if err := e.Load(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	e.WriteToc(rw)
}

func epubViewerHandler(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, "static/view.html")
}

func nextPageHandler(rw http.ResponseWriter, req *http.Request) {
	c, _ := req.Cookie("BookName")
	e := epub.New(strings.Replace(c.Value, "*", "\\", -1))
	if err := e.Load(); err != nil {
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
	c, _ := req.Cookie("BookName")
	e := epub.New(strings.Replace(c.Value, "*", "\\", -1))
	if err := e.Load(); err != nil {
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

	c, _ := req.Cookie("BookName")
	e := epub.New(strings.Replace(c.Value, "*", "\\", -1))
	if err := e.Load(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
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
	http.HandleFunc("/epubviewer/", epubViewerHandler)
	http.HandleFunc("/nextpage", nextPageHandler)
	http.HandleFunc("/prevpage", prevPageHandler)

	defer func() {
		err := recover()
		log.Println(err)
	}()

	http.ListenAndServe(":9090", nil)
}
