package service_test

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/service"
)

type MockCustomerRepo struct {
	mock.Mock
}

func (m *MockCustomerRepo) FindByID(id uint) (*models.Customer, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCustomerRepo) Create(customer *models.Customer) error {
	args := m.Called(customer)
	if args.Error(0) == nil {
		customer.ID = 1
	}
	return args.Error(0)
}

func (m *MockCustomerRepo) UpdateLoyalty(id uint, points int) error {
	args := m.Called(id, points)
	return args.Error(0)
}

func TestCustomerService_Register(t *testing.T) {
	repo := new(MockCustomerRepo)
	svc := service.NewCustomerService(repo)

	repo.On("Create", mock.AnythingOfType("*models.Customer")).Return(nil).Once()

	cust, err := svc.RegisterCustomer("Alice", "alice@test.com", "US")
	assert.NoError(t, err)
	assert.Equal(t, uint(1), cust.ID)
	assert.Equal(t, "New", cust.Segment)

	_, err2 := svc.RegisterCustomer("", "", "")
	assert.Error(t, err2)
	assert.Equal(t, "name and email cannot be empty", err2.Error())

	repo.On("Create", mock.AnythingOfType("*models.Customer")).Return(errors.New("db create error"))
	_, err3 := svc.RegisterCustomer("Fail", "fail@test.com", "US")
	assert.Error(t, err3)

	repo.AssertExpectations(t)
}

func TestCustomerService_Reward(t *testing.T) {
	repo := new(MockCustomerRepo)
	svc := service.NewCustomerService(repo)

	repo.On("UpdateLoyalty", uint(1), 50).Return(nil)

	err := svc.RewardCustomer(1, 50)
	assert.NoError(t, err)

	err2 := svc.RewardCustomer(1, -10)
	assert.Error(t, err2)
	assert.Equal(t, "reward points must be positive", err2.Error())

	repo.AssertExpectations(t)
}

func TestCustomerService_Get(t *testing.T) {
	repo := new(MockCustomerRepo)
	svc := service.NewCustomerService(repo)

	repo.On("FindByID", uint(1)).Return(&models.Customer{Name: "Bob"}, nil)

	cust, err := svc.GetCustomer(1)
	assert.NoError(t, err)
	assert.Equal(t, "Bob", cust.Name)

	_, err2 := svc.GetCustomer(0)
	assert.Error(t, err2)
	assert.Equal(t, "invalid id", err2.Error())

	// test repo error
	repo.On("FindByID", uint(99)).Return((*models.Customer)(nil), errors.New("not found"))
	_, err3 := svc.GetCustomer(99)
	assert.Error(t, err3)

	repo.AssertExpectations(t)
}
