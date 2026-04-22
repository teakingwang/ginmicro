package handler

import (
	"github.com/teakingwang/ginmicro/pkg/errs"
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
	req := &LoginReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		errs.ResponseError(c, errs.CodeInvalidArgs, err.Error())
		return
	}

	token, userDTO, err := h.svc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		errs.ResponseError(c, errs.CodeUnauthorized, err.Error())
		return
	}

	resp := &LoginResp{
		Token: token,
		User: &UserItem{
			ID:       userDTO.ID,
			Username: userDTO.Username,
			Nickname: userDTO.Nickname,
		},
	}

	errs.ResponseSuccessWithData(c, resp)
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

	errs.ResponseSuccessWithData(c, resp)
}
