package epub

import (
	"archive/zip"
	"encoding/xml"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
	"io"
	"strings"
)

type Ebook struct {
	fs      vfs.FileSystem
	opf     opf
	ncx     ncx
	TocPath string
	isOEBPS bool
}

func New(name string) (*Ebook, error) {
	rc, err := zip.OpenReader(name)
	if err != nil {
		return nil, err
	}
	//defer rc.Close()

	return &Ebook{fs: zipfs.New(rc, name)}, nil
}

func (e *Ebook) Load() error {
	//read .opf file
	if err := e.loadOpf(); err != nil {
		return err
	}

	//read .ncx file
	//if err := e.loadNcx(); err != nil {
	for _, item := range e.opf.Manifest.Item {
		if item.Id == e.opf.Spine.ItemRef[0].Idref {
			e.TocPath = item.Href
			break
		}
	}
	//}

	return nil
}

func (e *Ebook) loadOpf() error {
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

	err = xml.Unmarshal(buf, &e.opf)
	if err != nil {
		return err
	}

	return nil
}

func (e *Ebook) loadNcx() error {
	buf, err := vfs.ReadFile(e.fs, "/toc.ncx")
	if err != nil {
		return err
	}

	err = xml.Unmarshal(buf, &e.ncx)
	if err != nil {
		return err
	}

	return nil
}

func (e *Ebook) WriteFile(w io.Writer, path string) error {
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
