package service_test

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/service"
)

type MockStoreRepo struct {
	mock.Mock
}

func (m *MockStoreRepo) FindByID(id uint) (*models.Store, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Store), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStoreRepo) Create(store *models.Store) error {
	args := m.Called(store)
	if args.Error(0) == nil {
		store.ID = 1
	}
	return args.Error(0)
}

func (m *MockStoreRepo) GetStaffByStore(storeID uint) ([]models.Staff, error) {
	args := m.Called(storeID)
	return args.Get(0).([]models.Staff), args.Error(1)
}

func (m *MockStoreRepo) AddStaff(staff *models.Staff) error {
	args := m.Called(staff)
	if args.Error(0) == nil {
		staff.ID = 1
	}
	return args.Error(0)
}

func TestStoreService_CreateStore(t *testing.T) {
	repo := new(MockStoreRepo)
	svc := service.NewStoreService(repo)

	repo.On("Create", mock.AnythingOfType("*models.Store")).Return(nil).Once()

	store, err := svc.CreateStore("TestStore", "TestLoc", 0.1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), store.ID)

	_, err2 := svc.CreateStore("", "", 0)
	assert.Error(t, err2)

	repo.On("Create", mock.AnythingOfType("*models.Store")).Return(errors.New("db create error"))
	_, err3 := svc.CreateStore("FailStore", "Loc", 0.1)
	assert.Error(t, err3)

	repo.AssertExpectations(t)
}

func TestStoreService_HireStaff(t *testing.T) {
	repo := new(MockStoreRepo)
	svc := service.NewStoreService(repo)

	// Valid Store ID
	repo.On("FindByID", uint(1)).Return(&models.Store{ID: 1}, nil)
	repo.On("AddStaff", mock.AnythingOfType("*models.Staff")).Return(nil).Once()

	staff, err := svc.HireStaff(1, "Bob", "Clerk")
	assert.NoError(t, err)
	assert.Equal(t, uint(1), staff.ID)
	assert.Equal(t, uint(1), staff.StoreID)

	// DB error on AddStaff
	repo.On("AddStaff", mock.AnythingOfType("*models.Staff")).Return(errors.New("db insert error")).Once()
	_, errDb := svc.HireStaff(1, "FailBob", "Clerk")
	assert.Error(t, errDb)

	// missing params
	_, err2 := svc.HireStaff(1, "", "")
	assert.Error(t, err2)

	// Store not found
	repo.On("FindByID", uint(2)).Return((*models.Store)(nil), errors.New("db error"))
	_, err3 := svc.HireStaff(2, "Jim", "Mgr")
	assert.Error(t, err3)

	repo.AssertExpectations(t)
}

func TestStoreService_GetStoreDetails(t *testing.T) {
	repo := new(MockStoreRepo)
	svc := service.NewStoreService(repo)

	repo.On("FindByID", uint(1)).Return(&models.Store{ID: 1}, nil)
	repo.On("GetStaffByStore", uint(1)).Return([]models.Staff{{Name: "A"}, {Name: "B"}}, nil)

	store, staffs, err := svc.GetStoreDetails(1)
	assert.NoError(t, err)
	assert.NotNil(t, store)
	assert.Len(t, staffs, 2)

	repo.On("FindByID", uint(2)).Return((*models.Store)(nil), errors.New("db error"))
	_, _, err2 := svc.GetStoreDetails(2)
	assert.Error(t, err2)

	repo.AssertExpectations(t)
}
