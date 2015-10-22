package epub

import (
	"archive/zip"
)

var Cache map[string]*zip.ReadCloser

func init() {
	Cache = make(map[string]*zip.ReadCloser)
}
