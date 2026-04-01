package service

import (
	"errors"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
)

type AttributeService interface { AddAttribute(productID uint, key, value string) (*models.Attribute, error) }
type attrService struct{ repo repository.AttributeRepository }
func NewAttributeService(repo repository.AttributeRepository) AttributeService { return &attrService{repo} }
func (s *attrService) AddAttribute(productID uint, key, value string) (*models.Attribute, error) {
	if productID == 0 || key == "" { return nil, errors.New("invalid attribute") }
	attr := &models.Attribute{ProductID: productID, Key: key, Value: value}
	return attr, s.repo.Create(attr)
}
