package repository

import (
	"github.com/user/grafana-analytics-app/internal/models"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *models.Payment) error
	FindByOrderID(orderID uint) (*models.Payment, error)
}
type paymentRepo struct{ db *gorm.DB }
func NewPaymentRepository(db *gorm.DB) PaymentRepository { return &paymentRepo{db: db} }
func (r *paymentRepo) Create(payment *models.Payment) error { return r.db.Create(payment).Error }
func (r *paymentRepo) FindByOrderID(orderID uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}
