package service

import (
	"errors"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
	"time"
)

type PaymentService interface {
	ProcessPayment(orderID uint, method string) (*models.Payment, error)
}
type paymentService struct{ repo repository.PaymentRepository }
func NewPaymentService(repo repository.PaymentRepository) PaymentService { return &paymentService{repo} }
func (s *paymentService) ProcessPayment(orderID uint, method string) (*models.Payment, error) {
	if orderID == 0 || method == "" {
		return nil, errors.New("invalid payment info")
	}
	payment := &models.Payment{
		OrderID: orderID,
		Method:  method,
		Status:  "SUCCESS",
		PaidAt:  time.Now(),
	}
	err := s.repo.Create(payment)
	return payment, err
}
