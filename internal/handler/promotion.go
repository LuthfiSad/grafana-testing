package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/user/grafana-analytics-app/internal/service"
)

type PromotionHandler struct{ service service.PromotionService }
func NewPromotionHandler(svc service.PromotionService) *PromotionHandler { return &PromotionHandler{svc} }

type createPromoReq struct {
	Code     string  `json:"code" binding:"required"`
	Discount float64 `json:"discount" binding:"required"`
	Days     int     `json:"days"`
}

func (h *PromotionHandler) HandleCreatePromo(c *gin.Context) {
	var req createPromoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	promo, err := h.service.CreatePromo(req.Code, req.Discount, req.Days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"promotion": promo})
}
