package service_test

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/service"
)

type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) FindByID(id uint) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepo) Create(product *models.Product) error {
	args := m.Called(product)
	if args.Error(0) == nil {
		product.ID = 1
	}
	return args.Error(0)
}

func (m *MockProductRepo) UpdateStock(id uint, quantity int) error {
	args := m.Called(id, quantity)
	return args.Error(0)
}

func TestProductService_AddProduct(t *testing.T) {
	repo := new(MockProductRepo)
	svc := service.NewProductService(repo)

	repo.On("Create", mock.AnythingOfType("*models.Product")).Return(nil).Once()

	prod, err := svc.AddProduct("ItemA", "Cat1", 10.0, 5.0, 100)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), prod.ID)

	_, err2 := svc.AddProduct("", "Cat1", 10.0, 5.0, 100)
	assert.Error(t, err2)

	_, err3 := svc.AddProduct("ItemA", "Cat1", -10.0, 5.0, 100)
	assert.Error(t, err3)

	repo.On("Create", mock.AnythingOfType("*models.Product")).Return(errors.New("db create error"))
	_, err4 := svc.AddProduct("FailItem", "Cat1", 10.0, 5.0, 100)
	assert.Error(t, err4)

	repo.AssertExpectations(t)
}

func TestProductService_AdjustStock(t *testing.T) {
	repo := new(MockProductRepo)
	svc := service.NewProductService(repo)

	repo.On("UpdateStock", uint(1), 50).Return(nil)
	err := svc.AdjustStock(1, 50)
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestProductService_GetProduct(t *testing.T) {
	repo := new(MockProductRepo)
	svc := service.NewProductService(repo)

	repo.On("FindByID", uint(1)).Return(&models.Product{Name: "ItemA"}, nil)
	prod, err := svc.GetProduct(1)
	assert.NoError(t, err)
	assert.Equal(t, "ItemA", prod.Name)

	repo.On("FindByID", uint(2)).Return((*models.Product)(nil), errors.New("db error"))
	_, err2 := svc.GetProduct(2)
	assert.Error(t, err2)

	repo.AssertExpectations(t)
}
