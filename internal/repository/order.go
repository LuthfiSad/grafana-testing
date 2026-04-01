package repository

import (
	"github.com/user/grafana-analytics-app/internal/models"
	"gorm.io/gorm"
)

// OrderRepository is the interface for order data operations
type OrderRepository interface {
	InsertOrder(order *models.Order) error
	GetTotalRevenueByStatus(status string) (float64, error)
}

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new instance of OrderRepository
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// InsertOrder implements OrderRepository.
func (r *orderRepository) InsertOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

// GetTotalRevenueByStatus implements OrderRepository.
func (r *orderRepository) GetTotalRevenueByStatus(status string) (float64, error) {
	var total float64
	err := r.db.Model(&models.Order{}).
		Where("status = ?", status).
		Select("COALESCE(SUM(final_price), 0)").
		Scan(&total).Error
	return total, err
}
