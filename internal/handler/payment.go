package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/user/grafana-analytics-app/internal/service"
)

type PaymentHandler struct{ service service.PaymentService }
func NewPaymentHandler(svc service.PaymentService) *PaymentHandler { return &PaymentHandler{svc} }

type processPaymentReq struct {
	OrderID uint   `json:"order_id" binding:"required"`
	Method  string `json:"method" binding:"required"`
}

func (h *PaymentHandler) HandleProcessPayment(c *gin.Context) {
	var req processPaymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	payment, err := h.service.ProcessPayment(req.OrderID, req.Method)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"payment": payment})
}
