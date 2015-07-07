package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
	"html/template"
	"io"
	"strings"
)

type Epub struct {
	fs      vfs.FileSystem
	isOEBPS bool
	version int
}

func New(name string) (*Epub, error) {
	rc, err := zip.OpenReader(name)
	if err != nil {
		return nil, err
	}
	//defer rc.Close()

	return &Epub{fs: zipfs.New(rc, name)}, nil
}

func (e *Epub) WriteToc(w io.Writer) error {

	buf, err := vfs.ReadFile(e.fs, "/toc.ncx")
	if err != nil {
		return err
	}

	v := ncx{}

	err = xml.Unmarshal(buf, &v)
	if err != nil {
		return err
	}

	t, err := template.ParseFiles("ncx.template")
	if err != nil {
		return err
	}

	err = t.Execute(w, v)
	if err != nil {
		return err
	}

	return nil
}

func (e *Epub) WriteSpine(w io.Writer) error {

	buf, err := vfs.ReadFile(e.fs, "/content.opf")
	if strings.HasPrefix(err.Error(), "file not found") {
		buf, err = vfs.ReadFile(e.fs, "/OEBPS/content.opf")
		if err == nil {
			e.isOEBPS = true
		}
	}
	if err != nil {
		return err
	}

	v := OPF{}

	err = xml.Unmarshal(buf, &v)
	if err != nil {
		return err
	}

	for _, si := range v.Spine.ItemRefs {
		for _, mi := range v.Manifest.Items {
			if si.Idref == mi.Id {
				fmt.Fprintln(w, "<a href="+mi.Href+">"+si.Idref+"</a><br/>")
				break
			}
		}
	}

	return nil
}

func (e *Epub) WriteFile(w io.Writer, path string) error {
	if e.isOEBPS {
		path = "/OEBPS" + path
	}

	buf, err := vfs.ReadFile(e.fs, path)
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
