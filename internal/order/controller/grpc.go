package controller

import (
	"context"

	pb "github.com/teakingwang/ginmicro/api/order"
	"github.com/teakingwang/ginmicro/internal/order/service"
)

type OrderController struct {
	pb.UnimplementedOrderServiceServer
	svc service.OrderService
}

func NewOrderController(svc service.OrderService) pb.OrderServiceServer {
	return &OrderController{svc: svc}
}

func (oc *OrderController) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	dto, err := oc.svc.GetOrder(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if dto == nil {
		return nil, nil // 如果订单不存在，返回 nil
	}

	return &pb.GetOrderResponse{
		OrderID:  dto.OrderID,
		OrderSN:  dto.OrderSN,
		UserID:   dto.UserID,
		Username: dto.Username,
	}, nil
}
