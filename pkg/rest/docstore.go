package rest

import (
	"context"

	"github.com/sirupsen/logrus"
	_ "gocloud.dev/docstore/memdocstore"
	_ "gocloud.dev/docstore/mongodocstore"

	"gocloud.dev/docstore"
)

var Collection *docstore.Collection

// only call in main
func Connect() *docstore.Collection {
	if Collection == nil {
		var err error
		Collection, err = docstore.OpenCollection(context.Background(), mongoURL)
		if err != nil {
			logrus.Fatal(err)
		}
	}
	return Collection
}
