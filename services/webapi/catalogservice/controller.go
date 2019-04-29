package catalogservice

import (
	"github.com/gin-gonic/gin"
	"github.com/peizhong/letsgo/domain/catalog"
	"github.com/peizhong/letsgo/framework"
	log "github.com/sirupsen/logrus"
)

type ProductController struct {
	HTTPContext *gin.Context
	GoContext   *framework.GoContext
}

/*
每个函数最多只调一个domain的内容
*/

func (c ProductController) GetProduct() {
	log.Info("GetProduct")
	data, _ := catalog.GetProduct(c.GoContext, catalog.GetProductRequest{
		ProductId: c.HTTPContext.Param("id"),
	})
	c.HTTPContext.JSON(200, data)
}

func (c ProductController) GetProducts() {
	log.Info("GetProducts")
	pageIndex, _ := framework.IntTryParse(c.HTTPContext.Param("pageindex"))
	pageSize, _ := framework.IntTryParse(c.HTTPContext.DefaultQuery("pagesize", "100"))
	data, _ := catalog.GetProducts(c.GoContext, catalog.GetProductsRequest{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
	c.HTTPContext.JSON(200, data)
}
