package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/teakingwang/ginmicro/internal/order/handler"
	"github.com/teakingwang/ginmicro/internal/order/service"
)

func NewHTTPRouter(svc service.OrderService) *gin.Engine {
	r := gin.Default()
	h := handler.NewOrderHandler(svc)

	v1 := r.Group("/v1/order")
	{
		v1.GET("/:id", h.GetOrder) // 通过 :id 动态参数支持 REST 风格
	}

	return r
}
