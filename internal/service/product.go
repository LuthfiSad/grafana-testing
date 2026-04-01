package service

import (
	"errors"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
)

type ProductService interface {
	AddProduct(name, category string, price, cost float64, stock int) (*models.Product, error)
	AdjustStock(id uint, change int) error
	GetProduct(id uint) (*models.Product, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) AddProduct(name, category string, price, cost float64, stock int) (*models.Product, error) {
	if name == "" || price <= 0 {
		return nil, errors.New("invalid product details")
	}

	prod := &models.Product{
		Name:     name,
		Category: category,
		Price:    price,
		Cost:     cost,
		Stock:    stock,
	}

	err := s.repo.Create(prod)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func (s *productService) AdjustStock(id uint, change int) error {
	return s.repo.UpdateStock(id, change)
}

func (s *productService) GetProduct(id uint) (*models.Product, error) {
	return s.repo.FindByID(id)
}
