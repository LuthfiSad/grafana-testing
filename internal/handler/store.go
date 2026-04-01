package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/grafana-analytics-app/internal/service"
)

type StoreHandler struct {
	service service.StoreService
}

func NewStoreHandler(svc service.StoreService) *StoreHandler {
	return &StoreHandler{service: svc}
}

type createStoreReq struct {
	Name     string  `json:"name" binding:"required"`
	Location string  `json:"location" binding:"required"`
	TaxRate  float64 `json:"tax_rate"`
}

func (h *StoreHandler) HandleCreateStore(c *gin.Context) {
	var req createStoreReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	store, err := h.service.CreateStore(req.Name, req.Location, req.TaxRate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"store": store})
}
