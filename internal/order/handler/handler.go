package handler

import (
	"github.com/teakingwang/ginmicro/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/teakingwang/ginmicro/internal/order/service"
)

type OrderHandler struct {
	svc service.OrderService
}

func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	logger.Info("GetOrder called with ID:", id)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing order ID"})
		return
	}
	// 如果需要将 id 转换成 int 类型
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	resp, err := h.svc.GetOrder(c.Request.Context(), int64(idInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
