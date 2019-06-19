package catalog

import (
	"github.com/peizhong/letsgo/pkg/db"
	"github.com/peizhong/letsgo/playground/models/catalog"
)

type CatalogDomain struct {
}

func (*CatalogDomain) Gets() []*catalog.Classify {
	data := []*catalog.Classify{}
	db.Gets(&data)
	return data
}
