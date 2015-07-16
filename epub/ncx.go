package epub

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
