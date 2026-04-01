package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/grafana-analytics-app/internal/service"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{service: svc}
}

type addProductReq struct {
	Name     string  `json:"name" binding:"required"`
	Category string  `json:"category" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Cost     float64 `json:"cost" binding:"required"`
	Stock    int     `json:"stock"`
}

func (h *ProductHandler) HandleAddProduct(c *gin.Context) {
	var req addProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	prod, err := h.service.AddProduct(req.Name, req.Category, req.Price, req.Cost, req.Stock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"product": prod})
}
