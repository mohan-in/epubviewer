package epub

import (
	"archive/zip"
	"encoding/xml"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
)

func (e *Ebook) Load() error {

	r, err := zip.OpenReader(e.name)
	if err != nil {
		return err
	}
	defer r.Close()

	fs := zipfs.New(r, e.name)

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
