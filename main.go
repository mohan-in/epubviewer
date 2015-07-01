package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var e *epub

func uploadHandler(rw http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)

	src, _, err := r.FormFile("epubupload")
	if err != nil {
		fmt.Println(err)
	}

	dst, err := ioutil.TempFile("files", "file")
	if err != nil {
		fmt.Println(err)
	}
	defer dst.Close()

	io.Copy(dst, src)

	e, err = New(dst.Name())
	if err != nil {
		fmt.Println(err)
	}

	if err := e.WriteToc(rw); err != nil {
		if err := e.WriteSpine(rw); err != nil {
			fmt.Println(err)
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
