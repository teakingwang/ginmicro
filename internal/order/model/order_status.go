package model

type OrderStatus int

const (
	OrderStatusPending  OrderStatus = iota // 0：待支付
	OrderStatusPaid                        // 1：已支付
	OrderStatusCanceled                    // 2：已取消
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusPending:
		return "待支付"
	case OrderStatusPaid:
		return "已支付"
	case OrderStatusCanceled:
		return "已取消"
	default:
		return "未知状态"
	}
}
