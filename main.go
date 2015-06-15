package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
	"net/http"
)

var fs vfs.FileSystem

func init() {
	rc, err := zip.OpenReader("files/1. Looking for Alaska - John Green.epub")
	if err != nil {
		fmt.Println(err)
	}
	//defer rc.Close()

	fs = zipfs.New(rc, "Feynman")
}

func tocHandler(rw http.ResponseWriter, req *http.Request) {

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
		fmt.Fprintln(rw, "<a href="+t.Href+">"+t.Text+"</a>")
	}
}

func spineHandler(rw http.ResponseWriter, req *http.Request) {
	buf, err := vfs.ReadFile(fs, req.URL.Path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(rw, string(buf))
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

func main() {
	http.HandleFunc("/", spineHandler)
	http.HandleFunc("/toc", tocHandler)
	http.ListenAndServe(":9090", nil)
}
