package repository

import (
	"github.com/user/grafana-analytics-app/internal/models"
	"gorm.io/gorm"
)

type StoreRepository interface {
	FindByID(id uint) (*models.Store, error)
	Create(store *models.Store) error
	GetStaffByStore(storeID uint) ([]models.Staff, error)
	AddStaff(staff *models.Staff) error
}

type storeRepo struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepo{db: db}
}

func (r *storeRepo) FindByID(id uint) (*models.Store, error) {
	var store models.Store
	err := r.db.First(&store, id).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepo) Create(store *models.Store) error {
	return r.db.Create(store).Error
}

func (r *storeRepo) GetStaffByStore(storeID uint) ([]models.Staff, error) {
	var staff []models.Staff
	err := r.db.Where("store_id = ?", storeID).Find(&staff).Error
	return staff, err
}

func (r *storeRepo) AddStaff(staff *models.Staff) error {
	return r.db.Create(staff).Error
}
