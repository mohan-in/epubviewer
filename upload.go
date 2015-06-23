package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
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

	dir, err := ioutil.TempDir("c://wrk", "dir")
	if err != nil {
		fmt.Println(err)
	}

	dst, err := ioutil.TempFile(dir, "file")
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

	buf, err := vfs.ReadFile(fs, "/content.opf")
	if err != nil {
		fmt.Println(err)
	}

	v := OPF{}

	er := xml.Unmarshal(buf, &v)
	if err != nil {
		fmt.Println(er)
	}

	for _, t := range tocFromSpine(v.Spine, v.Manifest) {
		fmt.Fprintln(rw, "<a href="+t.Href+">"+t.Text+"</a><br/>")
	}
}

type toc struct {
	Text string
	Href string
}

func tocFromSpine(s spine, m manifest) []*toc {
	t := make([]*toc, 0)
	for _, si := range s.ItemRefs {
		tt := &toc{}
		for _, mi := range m.Items {
			if si.Idref == mi.Id {
				tt.Text = si.Idref
				tt.Href = mi.Href
				t = append(t, tt)
				break
			}
		}
	}
	return t
}

func spineHandler(rw http.ResponseWriter, req *http.Request) {
	if fs == nil {
		indexHandler(rw, req)
	}

	buf, err := vfs.ReadFile(fs, req.URL.Path)
	if err != nil {
		fmt.Println(err)
	}
	rw.Write(buf)
}

func indexHandler(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, "upload.html")
}

func main() {
	http.HandleFunc("/", spineHandler)
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/upload", uploadHandler)

	http.ListenAndServe(":9090", nil)
}
