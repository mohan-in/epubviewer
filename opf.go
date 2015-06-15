package main

import (
	"encoding/xml"
)

type OPF struct {
	XMLName  xml.Name `xml:"package"`
	Metadata metadata `xml:"metadata"`
	Manifest manifest `xml:"manifest"`
	Spine    spine    `xml:"spine"`
}

type metadata struct {
	Title   string   `xml:"title"`
	Creator []string `xml:"creator"`
	Lang    string   `xml:"language"`
}

type manifest struct {
	Items []item `xml:"item"`
}

type spine struct {
	TOC      string    `xml:"toc,attr"`
	ItemRefs []itemRef `xml:"itemref"`
}

type item struct {
	Id        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type"`
}

type itemRef struct {
	Idref string `xml:"idref,attr"`
}
