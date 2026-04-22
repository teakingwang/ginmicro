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

func (h *UserHandler) Login(c *gin.Context) {
	type LoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var req LoginRequest
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

func (h *UserHandler) GetUserList(c *gin.Context) {
	req := GetUserListReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		errs.ResponseError(c, errs.CodeInvalidArgs, err.Error())
		return
	}

	userDTOList, total, err := h.svc.GetUserList(c.Request.Context())
	if err != nil {
		errs.ResponseError(c, errs.CodeDatabaseError, err.Error())
		return
	}

	list := make([]*UserItem, 0, len(userDTOList))

	for _, dto := range userDTOList {
		list = append(list, &UserItem{
			ID:       dto.ID,
			Username: dto.Username,
			Nickname: dto.Nickname,
		})
	}

	resp := GetUserListResp{
		List:  list,
		Total: total,
	}

	errs.ResponseSuccessWithData(c, resp)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	logger.Info("GetUser called with ID:", id)
	if id == "" {
		errs.ResponseError(c, errs.CodeInvalidArgs, "missing user ID")
		return
	}
	// 如果需要将 id 转换成 int 类型
	idInt, err := strconv.Atoi(id)
	if err != nil {
		errs.ResponseError(c, errs.CodeInvalidArgs, "invalid user ID")
		return
	}

	userDTO, err := h.svc.GetUser(c.Request.Context(), int64(idInt))
	if err != nil {
		errs.ResponseError(c, errs.CodeDatabaseError, err.Error())
		return
	}

	resp := GetUserResp{
		User: &UserItem{
			ID:       userDTO.ID,
			Username: userDTO.Username,
			Nickname: userDTO.Nickname,
		},
	}

	c.JSON(http.StatusOK, resp)
}
