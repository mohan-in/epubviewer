package epub

import (
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
	"io"
)

func (e *Ebook) WriteToc(w io.Writer) error {
	return e.WriteFile(w, e.tocPath)
}

func (e *Ebook) WriteFile(w io.Writer, path string) error {

	r := Cache[e.Name]
	defer r.Close()

	fs := zipfs.New(r, e.Name)

	if e.isOEBPS {
		path = "/OEBPS" + path
	}

	buf, err := vfs.ReadFile(fs, path)
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
