package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
)

var fs vfs.FileSystem

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

	rc, err := zip.OpenReader(dst.Name())
	if err != nil {
		fmt.Println(err)
	}
	//defer rc.Close()

	fs = zipfs.New(rc, dst.Name())

	buf, err := vfs.ReadFile(fs, "/toc.ncx")
	if err == nil {
		v := ncx{}
		err := xml.Unmarshal(buf, &v)
		if err != nil {
			fmt.Println(err)
		}
		t, err := template.ParseFiles("ncx.template")
		if err != nil {
			fmt.Println(err)
		}
		t.Execute(rw, v)
		return
	} else {
		fmt.Println(err)
	}

	buf, err = vfs.ReadFile(fs, "/content.opf")
	if err != nil {
		fmt.Println(err)
	}

	v := OPF{}

	err = xml.Unmarshal(buf, &v)
	if err != nil {
		fmt.Println(err)
	}

	for _, t := range tocFromSpine(v.Spine, v.Manifest) {
		fmt.Fprintln(rw, "<a href="+t.Href+">"+t.Text+"</a><br/>")
	}
}

func spineHandler(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		http.ServeFile(rw, req, "index.html")
	} else {
		buf, err := vfs.ReadFile(fs, req.URL.Path)
		if err != nil {
			fmt.Println(err)
		}
		rw.Write(buf)
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
