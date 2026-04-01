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

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) AddProduct(n, c string, p, cost float64, s int) (*models.Product, error) {
	args := m.Called(n, c, p, cost, s)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductService) AdjustStock(id uint, change int) error {
	args := m.Called(id, change)
	return args.Error(0)
}

func (m *MockProductService) GetProduct(id uint) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestProductHandler_AddProduct(t *testing.T) {
	mockSvc := new(MockProductService)
	h := handler.NewProductHandler(mockSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/products", h.HandleAddProduct)

	mockSvc.On("AddProduct", "ItemA", "Cat1", 10.0, 5.0, 100).Return(&models.Product{ID: 1}, nil)

	body := []byte(`{"name":"ItemA","category":"Cat1","price":10.0,"cost":5.0,"stock":100}`)
	req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Test Invalid Payload
	bodyInv := []byte(`{}`)
	reqInv, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(bodyInv))
	reqInv.Header.Set("Content-Type", "application/json")
	wInv := httptest.NewRecorder()
	router.ServeHTTP(wInv, reqInv)
	assert.Equal(t, http.StatusBadRequest, wInv.Code)

	// Test Internal Error
	mockSvc.On("AddProduct", "ItemErr", "Cat1", 10.0, 5.0, 100).Return((*models.Product)(nil), errors.New("boom"))
	bodyErr := []byte(`{"name":"ItemErr","category":"Cat1","price":10.0,"cost":5.0,"stock":100}`)
	reqErr, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(bodyErr))
	reqErr.Header.Set("Content-Type", "application/json")
	wErr := httptest.NewRecorder()
	router.ServeHTTP(wErr, reqErr)
	assert.Equal(t, http.StatusInternalServerError, wErr.Code)
}
