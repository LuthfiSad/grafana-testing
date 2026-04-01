package service

import (
	"errors"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
)

type OrderItemService interface { AddItemToOrder(orderID, productID uint, qty int, price float64) (*models.OrderItem, error) }
type orderItemSvc struct{ repo repository.OrderItemRepository }
func NewOrderItemService(repo repository.OrderItemRepository) OrderItemService { return &orderItemSvc{repo} }
func (s *orderItemSvc) AddItemToOrder(orderID, productID uint, qty int, price float64) (*models.OrderItem, error) {
	if qty <= 0 { return nil, errors.New("invalid quantity") }
	item := &models.OrderItem{OrderID: orderID, ProductID: productID, Quantity: qty, UnitPrice: price}
	return item, s.repo.Create(item)
}
