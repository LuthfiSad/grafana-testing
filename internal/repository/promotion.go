package repository

import (
	"github.com/user/grafana-analytics-app/internal/models"
	"gorm.io/gorm"
)

type PromotionRepository interface {
	Create(promo *models.Promotion) error
	FindByCode(code string) (*models.Promotion, error)
}

type promotionRepo struct{ db *gorm.DB }

func NewPromotionRepository(db *gorm.DB) PromotionRepository { return &promotionRepo{db: db} }

func (r *promotionRepo) Create(promo *models.Promotion) error {
	return r.db.Create(promo).Error
}
func (r *promotionRepo) FindByCode(code string) (*models.Promotion, error) {
	var promo models.Promotion
	err := r.db.Where("code = ?", code).First(&promo).Error
	return &promo, err
}
