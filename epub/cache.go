package epub

import (
	"archive/zip"
)

type Cache map[string]*zip.ReadCloser

var cache Cache

func init() {
	cache = make(map[string]*zip.ReadCloser)
}
