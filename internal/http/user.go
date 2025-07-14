package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/teakingwang/ginmicro/internal/user/service"
	"net/http"
)

type UserHTTPController struct {
	svc service.UserService
}

func NewUserHTTPController(svc service.UserService) *UserHTTPController {
	return &UserHTTPController{svc: svc}
}

func (uc *UserHTTPController) SignUp(c *gin.Context) {
	// 解析请求参数并调用 uc.svc.SignUp
	c.JSON(http.StatusOK, gin.H{"msg": "signup successful"})
}

func (uc *UserHTTPController) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "login successful"})
}

func (uc *UserHTTPController) GetProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "profile info"})
}
