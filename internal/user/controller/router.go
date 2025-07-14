package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/teakingwang/ginmicro/internal/user/handler"
	"github.com/teakingwang/ginmicro/internal/user/service"
)

func NewHTTPRouter(svc service.UserService) *gin.Engine {
	r := gin.Default()

	h := handler.NewUserHandler(svc)

	v1 := r.Group("/v1/user")
	{
		v1.GET("/:id", h.GetUser)
	}

	return r
}
