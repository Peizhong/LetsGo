package main

import (
	"github.com/gin-gonic/gin"
	"github.com/peizhong/letsgo/domain/catalog"
	"github.com/peizhong/letsgo/framework"
	log "github.com/sirupsen/logrus"
)

type ProductController struct {
	HTTPContext *gin.Context
	DbContext   *framework.DbContext
}

func (c ProductController) GetProduct() {
	log.Info("GetProduct")
	data := catalog.GetProduct(c.DbContext, catalog.GetProductRequest{
		ProductId: c.HTTPContext.Param("id"),
	})
	c.HTTPContext.JSON(200, data)
}

func (c ProductController) GetProducts() {
	log.Info("GetProducts")
	pageIndex, _ := framework.IntTryParse(c.HTTPContext.Param("pageindex"))
	pageSize, _ := framework.IntTryParse(c.HTTPContext.DefaultQuery("pagesize", "100"))
	data := catalog.GetProducts(c.DbContext, catalog.GetProductsRequest{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
	c.HTTPContext.JSON(200, data)
}
