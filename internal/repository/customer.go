package repository

import (
	"github.com/user/grafana-analytics-app/internal/models"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	FindByID(id uint) (*models.Customer, error)
	Create(customer *models.Customer) error
	UpdateLoyalty(id uint, points int) error
}

type customerRepo struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepo{db: db}
}

func (r *customerRepo) FindByID(id uint) (*models.Customer, error) {
	var cust models.Customer
	err := r.db.First(&cust, id).Error
	if err != nil {
		return nil, err
	}
	return &cust, nil
}

func (r *customerRepo) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepo) UpdateLoyalty(id uint, points int) error {
	return r.db.Model(&models.Customer{}).Where("id = ?", id).Update("loyalty_points", gorm.Expr("loyalty_points + ?", points)).Error
}
