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
	type LoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误"})
		return
	}

	token, user, err := uc.svc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		},
	})
}

func (uc *UserHTTPController) GetProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "profile info"})
}