package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/teakingwang/ginmicro/internal/user/handler"
	"github.com/teakingwang/ginmicro/internal/user/service"
)

func NewHTTPRouter(svc service.UserService) *gin.Engine {
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	h := handler.NewUserHandler(svc)

	v1 := r.Group("/v1/user")
	{
		v1.POST("/login", h.Login)
		v1.GET("/list", h.GetUserList)
		v1.GET("/:id", h.GetUser)
	}

	return r
}
