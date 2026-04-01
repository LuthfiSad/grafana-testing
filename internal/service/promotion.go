package service

import (
	"errors"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
	"time"
)

type PromotionService interface {
	CreatePromo(code string, discount float64, daysValid int) (*models.Promotion, error)
	GetPromoByCode(code string) (*models.Promotion, error)
}

type promotionService struct {
	repo repository.PromotionRepository
}

func NewPromotionService(repo repository.PromotionRepository) PromotionService {
	return &promotionService{repo: repo}
}

func (s *promotionService) CreatePromo(code string, discount float64, daysValid int) (*models.Promotion, error) {
	if code == "" || discount <= 0 {
		return nil, errors.New("invalid promo data")
	}
	promo := &models.Promotion{
		Code:       code,
		Discount:   discount,
		ValidUntil: time.Now().AddDate(0, 0, daysValid),
	}
	err := s.repo.Create(promo)
	return promo, err
}

func (s *promotionService) GetPromoByCode(code string) (*models.Promotion, error) {
	return s.repo.FindByCode(code)
}
