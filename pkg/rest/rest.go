package rest

import "gocloud.dev/docstore"

type RestfulAPI struct {
	Coll *docstore.Collection
}

func NewRestfulAPI() *RestfulAPI {
	return &RestfulAPI{
		Coll: Connect(),
	}
}
