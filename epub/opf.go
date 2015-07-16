package epub

import (
	"encoding/xml"
)

type opf struct {
	XMLName xml.Name `xml:"package"`
	Version string   `xml:"version,attr"`

	Metadata struct {
		XMLName xml.Name `xml:"metadata"`
		Title   string   `xml:"title"`
		Creator []string `xml:"creator"`
		Lang    string   `xml:"language"`
	}

	Manifest struct {
		XMLName xml.Name `xml:"manifest"`

		Item []struct {
			XMLName   xml.Name `xml:"item"`
			Id        string   `xml:"id,attr"`
			Href      string   `xml:"href,attr"`
			MediaType string   `xml:"media-type"`
		}
	}

	Spine struct {
		XMLName xml.Name `xml:"spine"`
		TOC     string   `xml:"toc,attr"`

		ItemRef []struct {
			XMLName xml.Name `xml:"itemref"`
			Idref   string   `xml:"idref,attr"`
		}
	}

	Guide struct {
		XMLName   xml.Name `xml:"guide"`
		Reference []struct {
			XMLName xml.Name `xml:"reference"`
			toc     string   `xml:"title"`
			title   string   `xml:"title"`
			href    string   `xml:"href"`
		}
	}
}
