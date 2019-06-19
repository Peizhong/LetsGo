package catalog

import (
	"letsgo/framework/db"
	"letsgo/models/catalog"
)

type CatalogDomain struct {
}

func (*CatalogDomain) Gets() []*catalog.Classify {
	data := []*catalog.Classify{}
	db.Gets(&data)
	return data
}
