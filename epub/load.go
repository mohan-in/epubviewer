package epub

import (
	"archive/zip"
	"encoding/xml"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
	"io"
)

func (e *Ebook) Load(reader io.ReaderAt) error {

	var size int64
	for {
		p := make([]byte, 1)
		n, _ := reader.ReadAt(p, size)
		size = size + int64(n)
		if n == 0 {
			break
		}
	}

	r, _ := zip.NewReader(reader, size)

	rc := new(zip.ReadCloser)
	rc.Reader = *r

	fs := zipfs.New(rc, e.Name)

	//read .opf file
	if err := e.loadOpf(fs); err != nil {
		return err
	}

	//read .ncx file
	if err := e.loadNcx(fs); err != nil {
	}

	for _, item := range e.opf.Manifest.Item {
		if item.Id == e.opf.Spine.ItemRef[0].Idref {
			e.tocPath = "/" + item.Href
			break
		}
	}

	return nil
}

func (e *Ebook) loadOpf(fs vfs.FileSystem) error {
	buf, err := vfs.ReadFile(fs, "/content.opf")
	if err != nil {
		buf, err = vfs.ReadFile(fs, "/OEBPS/content.opf")
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

func (e *Ebook) loadNcx(fs vfs.FileSystem) error {
	buf, err := vfs.ReadFile(fs, "/toc.ncx")
	if err != nil {
		return err
	}

	err = xml.Unmarshal(buf, &e.ncx)
	if err != nil {
		return err
	}

	return nil
}
