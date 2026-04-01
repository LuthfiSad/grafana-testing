package service

import (
	"errors"
	"time"

	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
)

type CustomerService interface {
	RegisterCustomer(name, email, country string) (*models.Customer, error)
	RewardCustomer(id uint, points int) error
	GetCustomer(id uint) (*models.Customer, error)
}

type customerService struct {
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func (s *customerService) RegisterCustomer(name, email, country string) (*models.Customer, error) {
	if name == "" || email == "" {
		return nil, errors.New("name and email cannot be empty")
	}

	cust := &models.Customer{
		Name:          name,
		Email:         email,
		Country:       country,
		Segment:       "New",
		LoyaltyPoints: 0,
		CreatedAt:     time.Now(),
	}

	err := s.repo.Create(cust)
	if err != nil {
		return nil, err
	}
	return cust, nil
}

func (s *customerService) RewardCustomer(id uint, points int) error {
	if points <= 0 {
		return errors.New("reward points must be positive")
	}
	return s.repo.UpdateLoyalty(id, points)
}

func (s *customerService) GetCustomer(id uint) (*models.Customer, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	return s.repo.FindByID(id)
}
