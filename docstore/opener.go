package docstore

import (
	"context"
	"net/url"
)

// CollectionURLOpener opens a collection of documents based on a URL.
// The opener must not modify the URL argument. It must be safe to call from
// multiple goroutines.
//
// This interface is generally implemented by types in driver packages.
type CollectionURLOpener interface {
	OpenCollectionURL(ctx context.Context, u *url.URL) (*Docstore, error)
}
