package service

import (
	"errors"
	"time"

	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
)

// ProcessOrderRequest represents payload for creating an order
type ProcessOrderRequest struct {
	CustomerID    uint    `json:"customer_id"`
	StoreID       uint    `json:"store_id"`
	StaffReferral uint    `json:"staff_referral"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
}

// OrderService contains business logic methods
type OrderService interface {
	ProcessOrder(req ProcessOrderRequest) (*models.Order, error)
	CalculateTotalRevenue() (float64, error)
}

type orderService struct {
	repo repository.OrderRepository
}

// NewOrderService creates a new OrderService
func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

// ProcessOrder applies business logic before inserting order
func (s *orderService) ProcessOrder(req ProcessOrderRequest) (*models.Order, error) {
	if req.CustomerID == 0 || req.StoreID == 0 {
		return nil, errors.New("invalid customer or store ID")
	}
	if req.Amount <= 0 {
		return nil, errors.New("invalid order amount")
	}
	if req.Status == "" {
		req.Status = "PAID" // Default
	}

	order := &models.Order{
		CustomerID:    req.CustomerID,
		StoreID:       req.StoreID,
		StaffReferral: req.StaffReferral,
		SubTotal:      req.Amount,
		TaxAmount:     req.Amount * 0.1, // 10% tax in this region for example
		FinalPrice:    req.Amount * 1.1,
		Status:        req.Status,
		OrderDate:     time.Now(),
	}

	err := s.repo.InsertOrder(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// CalculateTotalRevenue implements OrderService.
func (s *orderService) CalculateTotalRevenue() (float64, error) {
	return s.repo.GetTotalRevenueByStatus("PAID")
}
