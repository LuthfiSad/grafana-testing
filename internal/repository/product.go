package repository

import (
	"github.com/user/grafana-analytics-app/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindByID(id uint) (*models.Product, error)
	Create(product *models.Product) error
	UpdateStock(id uint, quantity int) error
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Attributes").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepo) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepo) UpdateStock(id uint, quantity int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update("stock", gorm.Expr("stock + ?", quantity)).Error
}
