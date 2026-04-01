package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/user/grafana-analytics-app/internal/service"
)

type CustomerHandler struct {
	service service.CustomerService
}

func NewCustomerHandler(svc service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: svc}
}

type registerReq struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required"`
	Country string `json:"country"`
}

func (h *CustomerHandler) HandleRegister(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	cust, err := h.service.RegisterCustomer(req.Name, req.Email, req.Country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Customer registered", "customer": cust})
}

type rewardReq struct {
	Points int `json:"points" binding:"required"`
}

func (h *CustomerHandler) HandleReward(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var req rewardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid points payload"})
		return
	}

	err = h.service.RewardCustomer(uint(id), req.Points)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer rewarded"})
}
