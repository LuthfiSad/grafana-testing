package service_test

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/service"
)

// Mock Promo Repo
type MockPromoRepo struct{ mock.Mock }
func (m *MockPromoRepo) Create(p *models.Promotion) error { return m.Called(p).Error(0) }
func (m *MockPromoRepo) FindByCode(c string) (*models.Promotion, error) {
	args := m.Called(c)
	if args.Get(0) != nil { return args.Get(0).(*models.Promotion), args.Error(1) }
	return nil, args.Error(1)
}

func TestPromotionService(t *testing.T) {
	repo := new(MockPromoRepo)
	svc := service.NewPromotionService(repo)

	// Valid
	repo.On("Create", mock.Anything).Return(nil).Once()
	promo, err := svc.CreatePromo("SAVE10", 10.0, 30)
	assert.NoError(t, err)
	assert.NotNil(t, promo)

	// Invalid
	_, err2 := svc.CreatePromo("", 10.0, 30)
	assert.Error(t, err2)

	// Find
	repo.On("FindByCode", "SAVE10").Return(&models.Promotion{Code: "SAVE10"}, nil)
	found, err3 := svc.GetPromoByCode("SAVE10")
	assert.NoError(t, err3)
	assert.Equal(t, "SAVE10", found.Code)
}

// Mock Payment Repo
type MockPayRepo struct{ mock.Mock }
func (m *MockPayRepo) Create(p *models.Payment) error { return m.Called(p).Error(0) }
func (m *MockPayRepo) FindByOrderID(id uint) (*models.Payment, error) {
	args := m.Called(id)
	if args.Get(0) != nil { return args.Get(0).(*models.Payment), args.Error(1) }
	return nil, args.Error(1)
}

func TestPaymentService(t *testing.T) {
	repo := new(MockPayRepo)
	svc := service.NewPaymentService(repo)
	
	repo.On("Create", mock.Anything).Return(nil).Once()
	pay, err := svc.ProcessPayment(1, "CC")
	assert.NoError(t, err)
	assert.Equal(t, "CC", pay.Method)

	_, err2 := svc.ProcessPayment(0, "CC")
	assert.Error(t, err2)
}

// Mock Shared Repo for generic create
type MockGenericRepo struct{ mock.Mock }
func (m *MockGenericRepo) Create(entity interface{}) error { return m.Called(entity).Error(0) }

type MockShippingRepo struct { MockGenericRepo }
func (m *MockShippingRepo) Create(s *models.Shipping) error { return m.Called(s).Error(0) }
type MockReviewRepo struct { MockGenericRepo }
func (m *MockReviewRepo) Create(r *models.Review) error { return m.Called(r).Error(0) }
type MockRefundRepo struct { MockGenericRepo }
func (m *MockRefundRepo) Create(r *models.Refund) error { return m.Called(r).Error(0) }
type MockInvLogRepo struct { MockGenericRepo }
func (m *MockInvLogRepo) Create(i *models.InventoryLog) error { return m.Called(i).Error(0) }
type MockAttrRepo struct { MockGenericRepo }
func (m *MockAttrRepo) Create(a *models.Attribute) error { return m.Called(a).Error(0) }
type MockOrderItemRepo struct { MockGenericRepo }
func (m *MockOrderItemRepo) Create(i *models.OrderItem) error { return m.Called(i).Error(0) }

func TestAdditionalServices(t *testing.T) {
	// Shipping
	shipRepo := new(MockShippingRepo)
	shipRepo.On("Create", mock.Anything).Return(nil).Once()
	shipSvc := service.NewShippingService(shipRepo)
	_, err := shipSvc.ArrangeShipping(1, "FedEx", 15.0)
	assert.NoError(t, err)
	_, err = shipSvc.ArrangeShipping(0, "", 0.0)
	assert.Error(t, err)

	// Review
	revRepo := new(MockReviewRepo)
	revRepo.On("Create", mock.Anything).Return(nil).Once()
	revSvc := service.NewReviewService(revRepo)
	_, err = revSvc.LeaveReview(1, 1, 5, "Great")
	assert.NoError(t, err)
	_, err = revSvc.LeaveReview(1, 1, 6, "Bad Rating") // Invalid rating
	assert.Error(t, err)

	// Refund
	refRepo := new(MockRefundRepo)
	refRepo.On("Create", mock.Anything).Return(nil).Once()
	refSvc := service.NewRefundService(refRepo)
	_, err = refSvc.ProcessRefund(1, 10.0, "Broken")
	assert.NoError(t, err)
	_, err = refSvc.ProcessRefund(0, 10.0, "")
	assert.Error(t, err)

	// InvLog
	invRepo := new(MockInvLogRepo)
	invRepo.On("Create", mock.Anything).Return(nil).Once()
	invSvc := service.NewInventoryLogService(invRepo)
	_, err = invSvc.AddLog(1, 10, "Restock")
	assert.NoError(t, err)
	_, err = invSvc.AddLog(0, 10, "")
	assert.Error(t, err)

	// Attribute
	attrRepo := new(MockAttrRepo)
	attrRepo.On("Create", mock.Anything).Return(nil).Once()
	attrSvc := service.NewAttributeService(attrRepo)
	_, err = attrSvc.AddAttribute(1, "Color", "Red")
	assert.NoError(t, err)
	_, err = attrSvc.AddAttribute(0, "", "")
	assert.Error(t, err)

	// OrderItem
	oiRepo := new(MockOrderItemRepo)
	oiRepo.On("Create", mock.Anything).Return(nil).Once()
	oiSvc := service.NewOrderItemService(oiRepo)
	_, err = oiSvc.AddItemToOrder(1, 1, 2, 10.0)
	assert.NoError(t, err)
	_, err = oiSvc.AddItemToOrder(1, 1, 0, 10.0)
	assert.Error(t, err)
	
	// Add assertions for repo Create errors to get 100% coverage
	shipRepo.On("Create", mock.Anything).Return(errors.New("db error")).Once()
	_, err = shipSvc.ArrangeShipping(1, "FedEx", 15.0)
	assert.Error(t, err)

	revRepo.On("Create", mock.Anything).Return(errors.New("db error")).Once()
	_, err = revSvc.LeaveReview(1, 1, 5, "Great")
	assert.Error(t, err)

	refRepo.On("Create", mock.Anything).Return(errors.New("db error")).Once()
	_, err = refSvc.ProcessRefund(1, 10.0, "Broken")
	assert.Error(t, err)

	invRepo.On("Create", mock.Anything).Return(errors.New("db error")).Once()
	_, err = invSvc.AddLog(1, 10, "Restock")
	assert.Error(t, err)

	attrRepo.On("Create", mock.Anything).Return(errors.New("db error")).Once()
	_, err = attrSvc.AddAttribute(1, "Color", "Red")
	assert.Error(t, err)

	oiRepo.On("Create", mock.Anything).Return(errors.New("db error")).Once()
	_, err = oiSvc.AddItemToOrder(1, 1, 2, 10.0)
	assert.Error(t, err)
}
