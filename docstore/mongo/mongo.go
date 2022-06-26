package mongo

import (
	"context"
	"net/url"

	"github.com/silverswords/devcloud/docstore"
)

const Scheme = "mongo"

func init() {
	docstore.DefaultDriverMap().Register(Scheme, new(Opener))
}

type docStore struct {
}

func (d *docStore) Exec() {

}

func newDocStore() *docStore {
	return &docStore{}
}

func openDocstore() (*docstore.Docstore, error) {
	dc := newDocStore()
	return docstore.NewDocstore(dc), nil
}

type Opener struct {
}

func (o *Opener) OpenCollectionURL(ctx context.Context, u *url.URL) (*docstore.Docstore, error) {
	panic("not implement")
	return nil, nil
}
