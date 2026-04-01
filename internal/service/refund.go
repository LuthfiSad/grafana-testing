package service

import (
	"errors"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
	"time"
)

type RefundService interface { ProcessRefund(orderID uint, amount float64, reason string) (*models.Refund, error) }
type refundService struct{ repo repository.RefundRepository }
func NewRefundService(repo repository.RefundRepository) RefundService { return &refundService{repo} }
func (s *refundService) ProcessRefund(orderID uint, amount float64, reason string) (*models.Refund, error) {
	if orderID == 0 || amount <= 0 { return nil, errors.New("invalid refund data") }
	ref := &models.Refund{OrderID: orderID, Amount: amount, Reason: reason, CreatedAt: time.Now()}
	return ref, s.repo.Create(ref)
}
