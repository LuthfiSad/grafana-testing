package service

import (
	"errors"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
)

type ShippingService interface { ArrangeShipping(orderID uint, carrier string, cost float64) (*models.Shipping, error) }
type shippingService struct{ repo repository.ShippingRepository }
func NewShippingService(repo repository.ShippingRepository) ShippingService { return &shippingService{repo} }
func (s *shippingService) ArrangeShipping(orderID uint, carrier string, cost float64) (*models.Shipping, error) {
	if orderID == 0 || carrier == "" { return nil, errors.New("invalid shipping") }
	sh := &models.Shipping{OrderID: orderID, Carrier: carrier, ShippingCost: cost, EstimatedDays: 3}
	err := s.repo.Create(sh)
	return sh, err
}
