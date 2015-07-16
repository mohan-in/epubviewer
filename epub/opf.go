package epub

import (
	"encoding/xml"
)

type opf struct {
	XMLName  xml.Name `xml:"package"`
	Version  string   `xml:"version,attr"`
	Metadata metadata `xml:"metadata"`
	Manifest manifest `xml:"manifest"`
	Spine    spine    `xml:"spine"`
	Guide    guide    `xml:"guide"`
}

type metadata struct {
	Title   string   `xml:"title"`
	Creator []string `xml:"creator"`
	Lang    string   `xml:"language"`
}

type manifest struct {
	Item []item `xml:"item"`
}

type item struct {
	Id        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type"`
}

type spine struct {
	TOC     string    `xml:"toc,attr"`
	ItemRef []itemRef `xml:"itemref"`
}

type itemRef struct {
	Idref string `xml:"idref,attr"`
}

type guide struct {
	Reference []reference `xml:"reference"`
}

type reference struct {
	toc   string `xml:"title"`
	title string `xml:"title"`
	href  string `xml:"href"`
}
