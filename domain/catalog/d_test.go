package catalog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCatalogDomain_Gets(t *testing.T) {
	data := (&CatalogDomain{}).Gets()
	assert.NotEmpty(t, data)
}
