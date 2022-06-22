package mongo

import (
	"github.com/silverswords/devcloud/docstore"
)

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
