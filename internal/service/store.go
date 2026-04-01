package service

import (
	"errors"

	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
)

type StoreService interface {
	CreateStore(name, location string, taxRate float64) (*models.Store, error)
	HireStaff(storeID uint, name, role string) (*models.Staff, error)
	GetStoreDetails(storeID uint) (*models.Store, []models.Staff, error)
}

type storeService struct {
	repo repository.StoreRepository
}

func NewStoreService(repo repository.StoreRepository) StoreService {
	return &storeService{repo: repo}
}

func (s *storeService) CreateStore(name, location string, taxRate float64) (*models.Store, error) {
	if name == "" || location == "" {
		return nil, errors.New("name and location are required")
	}

	store := &models.Store{
		Name:     name,
		Location: location,
		TaxRate:  taxRate,
	}

	err := s.repo.Create(store)
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (s *storeService) HireStaff(storeID uint, name, role string) (*models.Staff, error) {
	if name == "" || role == "" {
		return nil, errors.New("staff name and role are required")
	}
	// Verify store exists
	_, err := s.repo.FindByID(storeID)
	if err != nil {
		return nil, errors.New("store not found")
	}

	staff := &models.Staff{
		StoreID: storeID,
		Name:    name,
		Role:    role,
	}

	err = s.repo.AddStaff(staff)
	if err != nil {
		return nil, err
	}
	return staff, nil
}

func (s *storeService) GetStoreDetails(storeID uint) (*models.Store, []models.Staff, error) {
	store, err := s.repo.FindByID(storeID)
	if err != nil {
		return nil, nil, errors.New("store not found")
	}

	staff, _ := s.repo.GetStaffByStore(storeID)
	return store, staff, nil
}
