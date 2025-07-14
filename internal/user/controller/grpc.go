package controller

import (
	"context"
	"github.com/teakingwang/ginmicro/pkg/logger"

	"github.com/teakingwang/ginmicro/api/user"
	"github.com/teakingwang/ginmicro/internal/user/service"
)

// userController 实现 user.UserServiceServer（由 proto 自动生成）
type userController struct {
	svc service.UserService
	// 必须嵌入这个匿名字段，否则会提示 missing mustEmbedUnimplementedUserServiceServer
	user.UnimplementedUserServiceServer
}

// NewUserController 构造函数
func NewUserController(svc service.UserService) user.UserServiceServer {
	return &userController{
		svc: svc,
	}
}

// 实现 gRPC 的 GetUser 方法
func (uc *userController) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	logger.Info("GetUser called with ID:", req.Id)
	// 调用 service 处理逻辑
	dto, err := uc.svc.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if dto == nil {
		return nil, nil // 如果用户不存在，返回 nil
	}

	return &user.GetUserResponse{
		Id:       dto.ID,
		Username: dto.Username,
		Email:    dto.Email,
	}, nil
}
