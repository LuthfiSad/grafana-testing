package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/grafana-analytics-app/internal/handler"
	"github.com/user/grafana-analytics-app/internal/models"
)

type MockCustomerService struct {
	mock.Mock
}

func (m *MockCustomerService) RegisterCustomer(name, email, country string) (*models.Customer, error) {
	args := m.Called(name, email, country)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCustomerService) RewardCustomer(id uint, points int) error {
	args := m.Called(id, points)
	return args.Error(0)
}

func (m *MockCustomerService) GetCustomer(id uint) (*models.Customer, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupCustomerRouter(mockSvc *MockCustomerService) *gin.Engine {
	custHandler := handler.NewCustomerHandler(mockSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/customers", custHandler.HandleRegister)
	router.POST("/api/customers/:id/reward", custHandler.HandleReward)
	return router
}

func TestCustomerHandler_Register(t *testing.T) {
	mockSvc := new(MockCustomerService)
	router := setupCustomerRouter(mockSvc)

	mockSvc.On("RegisterCustomer", "Tom", "tom@mail.com", "UK").Return(&models.Customer{ID: 1, Name: "Tom"}, nil)

	body := []byte(`{"name":"Tom","email":"tom@mail.com","country":"UK"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	// Test Invalid JSON
	reqBad, _ := http.NewRequest(http.MethodPost, "/api/customers", bytes.NewBuffer([]byte(`{bad}`)))
	wBad := httptest.NewRecorder()
	router.ServeHTTP(wBad, reqBad)
	assert.Equal(t, http.StatusBadRequest, wBad.Code)
}

func TestCustomerHandler_Register_ServiceErr(t *testing.T) {
	mockSvc := new(MockCustomerService)
	router := setupCustomerRouter(mockSvc)

	mockSvc.On("RegisterCustomer", "Sam", "sam@mail.com", "US").Return((*models.Customer)(nil), errors.New("db err"))

	body := []byte(`{"name":"Sam","email":"sam@mail.com","country":"US"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCustomerHandler_Reward(t *testing.T) {
	mockSvc := new(MockCustomerService)
	router := setupCustomerRouter(mockSvc)

	mockSvc.On("RewardCustomer", uint(1), 100).Return(nil)

	body := []byte(`{"points":100}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/customers/1/reward", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Test bad ID
	reqBadID, _ := http.NewRequest(http.MethodPost, "/api/customers/abc/reward", bytes.NewBuffer(body))
	wBadID := httptest.NewRecorder()
	router.ServeHTTP(wBadID, reqBadID)
	assert.Equal(t, http.StatusBadRequest, wBadID.Code)
	
	// Test bad points
	reqBadPoints, _ := http.NewRequest(http.MethodPost, "/api/customers/1/reward", bytes.NewBuffer([]byte(`{}`)))
	wBadPoints := httptest.NewRecorder()
	router.ServeHTTP(wBadPoints, reqBadPoints)
	assert.Equal(t, http.StatusBadRequest, wBadPoints.Code)
}

func TestCustomerHandler_Reward_ServiceErr(t *testing.T) {
	mockSvc := new(MockCustomerService)
	router := setupCustomerRouter(mockSvc)

	mockSvc.On("RewardCustomer", uint(1), 100).Return(errors.New("db error"))

	body := []byte(`{"points":100}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/customers/1/reward", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
