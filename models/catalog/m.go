package catalog

import (
	"time"

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
	Id          string    `gorm:"primary_key;column:Id"`
	WorkspaceId string    `gorm:"column:WorkspaceId"`
	Name        string    `gorm:"column:Name"`
	Code        string    `gorm:"column:Code"`
	ClassifyId  uint      `gorm:"column:CategoryId"`
	AssetState  uint      `gorm:"column:DataStatus"`
	Manufacture string    `gorm:"column:Manufacture"`
	UpdateTime  time.Time `gorm:"column:UpdateAt"`
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
	Id          string
	WorkspaceId string
	Name        string
	Code        string
	ClassifyId  uint
	AssetState  uint
	Manufacture string
}

type GetProductsResponse struct {
	Items     []*GetProductResponse
	Count     int
	PageIndex int
	PageSize  int
}
