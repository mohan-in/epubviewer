package epub

type Ebook struct {
	Name    string
	opf     opf
	ncx     ncx
	tocPath string
	isOEBPS bool
}

func New(name string) *Ebook {
	return &Ebook{Name: name}
}

func (e *Ebook) GetToc() string {
	return e.tocPath
}

func (e *Ebook) getManifestItemFromHref(href string) string {
	for _, item := range e.opf.Manifest.Item {
		if item.Href == href {
			return item.Id
		}
	}

	return ""
}

func (e *Ebook) getManifestItemFromId(id string) string {
	for _, item := range e.opf.Manifest.Item {
		if item.Id == id {
			return item.Href
		}
	}

	return ""
}

func (e *Ebook) GetNextPage(href string) string {

	id := e.getManifestItemFromHref(href)

	var nextId string
	for i, itemRef := range e.opf.Spine.ItemRef {
		if itemRef.Idref == id {
			nextId = e.opf.Spine.ItemRef[i+1].Idref
			break
		}
	}

	return e.getManifestItemFromId(nextId)
}

func (e *Ebook) GetPrevPage(href string) string {
	id := e.getManifestItemFromHref(href)

	var nextId string
	for i, itemRef := range e.opf.Spine.ItemRef {
		if itemRef.Idref == id {
			nextId = e.opf.Spine.ItemRef[i-1].Idref
			break
		}
	}

	return e.getManifestItemFromId(nextId)
}
