package repository

import (
	"github.com/user/grafana-analytics-app/internal/models"
	"gorm.io/gorm"
)

type OrderItemRepository interface { Create(item *models.OrderItem) error }
type orderItemRepo struct{ db *gorm.DB }
func NewOrderItemRepository(db *gorm.DB) OrderItemRepository { return &orderItemRepo{db: db} }
func (r *orderItemRepo) Create(item *models.OrderItem) error { return r.db.Create(item).Error }
