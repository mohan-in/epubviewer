package epub

import (
	"archive/zip"
	"bytes"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
	"io"
	"io/ioutil"
	"net/http"
)

type Cache map[string]*zip.ReadCloser

var (
	cache  Cache
	bucket string = "mohanram-epubreader.appspot.com"
)

func init() {
	cache = make(map[string]*zip.ReadCloser)
}

func AddToCache(c context.Context, filename string, src io.Reader) {
	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(c, storage.ScopeFullControl),
		},
	}
	ctx := cloud.NewContext(appengine.AppID(c), hc)
	dst := storage.NewWriter(ctx, bucket, filename)
	defer dst.Close()
	io.Copy(dst, src)
}

func GetFromCache(c context.Context, filename string) io.ReaderAt {
	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(c, storage.ScopeFullControl),
		},
	}
	ctx := cloud.NewContext(appengine.AppID(c), hc)
	r, _ := storage.NewReader(ctx, bucket, filename)
	defer r.Close()
	buf, _ := ioutil.ReadAll(r)
	return bytes.NewReader(buf)
}
