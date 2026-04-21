package handler

import (
	"github.com/teakingwang/ginmicro/pkg/errs"
	"net/http"
	"strconv"

	"github.com/teakingwang/ginmicro/api/user"
	"github.com/teakingwang/ginmicro/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/teakingwang/ginmicro/internal/user/service"
)

type UserHandler struct {
	svc service.UserService
	user.UnimplementedUserServiceServer
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// GET /v1/user/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	logger.Info("GetUser called with ID:", id)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user ID"})
		return
	}
	// 如果需要将 id 转换成 int 类型
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	resp, err := h.svc.GetUser(c.Request.Context(), int64(idInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetUserList(c *gin.Context) {
	req := GetUserListReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误"})
		return
	}

	userList, total, err := h.svc.GetUserList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": errs.CodeSuccess,
		"msg":  "success",
		"data": gin.H{
			"list":  userList,
			"total": total,
		},
	})
}

// POST /v1/user/signin
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误"})
		return
	}

	token, user, err := h.svc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": errs.CodeSuccess,
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
