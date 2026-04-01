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

type MockStoreService struct {
	mock.Mock
}

func (m *MockStoreService) CreateStore(n, l string, tr float64) (*models.Store, error) {
	args := m.Called(n, l, tr)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Store), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStoreService) HireStaff(sid uint, n, r string) (*models.Staff, error) {
	args := m.Called(sid, n, r)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Staff), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStoreService) GetStoreDetails(sid uint) (*models.Store, []models.Staff, error) {
	args := m.Called(sid)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Store), args.Get(1).([]models.Staff), args.Error(2)
	}
	return nil, nil, args.Error(2)
}

func TestStoreHandler_Create(t *testing.T) {
	mockSvc := new(MockStoreService)
	h := handler.NewStoreHandler(mockSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/stores", h.HandleCreateStore)

	mockSvc.On("CreateStore", "StoreA", "NYC", 0.1).Return(&models.Store{ID: 1}, nil)

	body := []byte(`{"name":"StoreA","location":"NYC","tax_rate":0.1}`)
	req, _ := http.NewRequest(http.MethodPost, "/stores", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Test Internal Error
	mockSvc.On("CreateStore", "StoreErr", "NYC", 0.1).Return((*models.Store)(nil), errors.New("boom"))
	bodyErr := []byte(`{"name":"StoreErr","location":"NYC","tax_rate":0.1}`)
	reqErr, _ := http.NewRequest(http.MethodPost, "/stores", bytes.NewBuffer(bodyErr))
	reqErr.Header.Set("Content-Type", "application/json")

	wErr := httptest.NewRecorder()
	router.ServeHTTP(wErr, reqErr)

	assert.Equal(t, http.StatusInternalServerError, wErr.Code)

	// Test invalid payload
	bodyInv := []byte(`{}`)
	reqInv, _ := http.NewRequest(http.MethodPost, "/stores", bytes.NewBuffer(bodyInv))
	reqInv.Header.Set("Content-Type", "application/json")

	wInv := httptest.NewRecorder()
	router.ServeHTTP(wInv, reqInv)

	assert.Equal(t, http.StatusBadRequest, wInv.Code)
}
