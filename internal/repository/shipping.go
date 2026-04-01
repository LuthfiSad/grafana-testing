package repository

import (
	"github.com/user/grafana-analytics-app/internal/models"
	"gorm.io/gorm"
)

type ShippingRepository interface { Create(shipping *models.Shipping) error }
type shippingRepo struct{ db *gorm.DB }
func NewShippingRepository(db *gorm.DB) ShippingRepository { return &shippingRepo{db: db} }
func (r *shippingRepo) Create(shipping *models.Shipping) error { return r.db.Create(shipping).Error }
