package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/user/grafana-analytics-app/internal/service"
)

// OrderHandler struct
type OrderHandler struct {
	service         service.OrderService
	ordersProcessed prometheus.Counter
	totalRevenue    *prometheus.CounterVec
}

// NewOrderHandler initializes OrderHandler
func NewOrderHandler(svc service.OrderService, processed prometheus.Counter, revenue *prometheus.CounterVec) *OrderHandler {
	return &OrderHandler{
		service:         svc,
		ordersProcessed: processed,
		totalRevenue:    revenue,
	}
}

// HandleProcessOrder is a gin handler that takes an order request, delegates to service, and records metrics
func (h *OrderHandler) HandleProcessOrder(c *gin.Context) {
	var req service.ProcessOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	order, err := h.service.ProcessOrder(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Metrics update
	h.ordersProcessed.Inc()

	// Logic from previous simple implementaion
	countries := []string{"Indonesia", "USA", "Germany", "Japan"}
	country := countries[order.CustomerID%uint(len(countries))]
	h.totalRevenue.WithLabelValues(country).Add(order.FinalPrice)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order processed successfully",
		"order":   order,
	})
}
