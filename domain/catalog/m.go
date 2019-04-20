package catalog

import (
	"github.com/jinzhu/gorm"
)

/*
model for db
*/
type BaseInfo struct {
	gorm.Model
	BaseInfoTypeId uint
	Name           string
	Column         string
	SortNo         uint
}

type Classify struct {
	Id             string `gorm:"primary_key;column:Id"`
	Name           string `gorm:"column:Name"`
	FullPath       string `gorm:"column:FullPath"`
	BaseInfoTypeId string `gorm:"_"`
	ParentId       uint   `gorm:"column:ParentId"`
}

func (Classify) TableName() string {
	return "Classifies"
}

type Product struct {
	Id          string
	WorkspaceId string
	Name        string
	Code        string
	ClassifyId  uint
	AssetState  uint
	Manufacture string
}

func (Product) TableName() string {
	return "Products"
}

/*
request
*/

type GetClassifyRequest struct {
	ClassifyId string
}

type GetClassifiesRequest struct {
}

type GetProductRequest struct {
	ProductId string
}

type GetProductsRequest struct {
	PageIndex int
	PageSize  int
}

/*
response
*/

type GetClassifyResponse struct {
}

type GetProductResponse struct {
}

type GetProductsResponse struct {
	Count     int
	PageIndex int
	PageSize  int
}
