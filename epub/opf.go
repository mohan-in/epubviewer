package epub

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

type ncx struct {
	NavMap navMap `xml:"navMap"`
}

type navMap struct {
	NavPoint []navPoint `xml:"navPoint"`
}

type navPoint struct {
	NavLabel navLabel `xml:"navLabel"`
	Content  content  `xml:"content"`
}

type navLabel struct {
	Text string `xml:"text"`
}

type content struct {
	Src string `xml:"src,attr"`
}

type toc struct {
	Text string
	Href string
}

func tocFromSpine(s spine, m manifest) []*toc {
	t := make([]*toc, 0)
	for _, si := range s.ItemRefs {
		tt := &toc{}
		for _, mi := range m.Items {
			if si.Idref == mi.Id {
				tt.Text = si.Idref
				tt.Href = mi.Href
				t = append(t, tt)
				break
			}
		}
	}
	return t
}
