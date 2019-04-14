package catalog

import (
	"github.com/gomodule/redigo/redis"
	"github.com/peizhong/letsgo/framework"
)

func GetProduct(c *framework.DbContext, r GetProductRequest) GetProductsResponse {
	var classify Classify
	c.Database.First(&classify, r.ProductId)
	// play redis
	_, err := c.Cache.Do("SET", "go_key", "redigo")
	if err == nil {
		v, err := redis.String(c.Cache.Do("GET", "go_key"))
		if err == nil {
			_ = v
		}
	}
	return GetProductsResponse{
		PageIndex: 1,
		PageSize:  100,
	}
}

func GetProducts(c *framework.DbContext, r GetProductsRequest) GetProductsResponse {
	var classifies []Classify
	c.Database.Offset(r.PageIndex * r.PageSize).Limit(r.PageSize).Find(&classifies)
	return GetProductsResponse{
		PageIndex: r.PageIndex,
		PageSize:  r.PageSize,
	}
}
