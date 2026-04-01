package repository

import (
	"github.com/user/grafana-analytics-app/internal/models"
	"gorm.io/gorm"
)

type RefundRepository interface { Create(refund *models.Refund) error }
type refundRepo struct{ db *gorm.DB }
func NewRefundRepository(db *gorm.DB) RefundRepository { return &refundRepo{db: db} }
func (r *refundRepo) Create(refund *models.Refund) error { return r.db.Create(refund).Error }
