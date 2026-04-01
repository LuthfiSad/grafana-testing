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

type MockPromotionSvc struct{ mock.Mock }
func (m *MockPromotionSvc) CreatePromo(c string, d float64, days int) (*models.Promotion, error) {
	args := m.Called(c, d, days)
	if args.Get(0) != nil { return args.Get(0).(*models.Promotion), args.Error(1) }
	return nil, args.Error(1)
}
func (m *MockPromotionSvc) GetPromoByCode(c string) (*models.Promotion, error) {
	args := m.Called(c)
	if args.Get(0) != nil { return args.Get(0).(*models.Promotion), args.Error(1) }
	return nil, args.Error(1)
}

func TestPromotionHandler_CreatePromo(t *testing.T) {
	mockSvc := new(MockPromotionSvc)
	h := handler.NewPromotionHandler(mockSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/promotions", h.HandleCreatePromo)

	mockSvc.On("CreatePromo", "SALE", 20.0, 10).Return(&models.Promotion{Code: "SALE"}, nil)

	body := []byte(`{"code":"SALE","discount":20.0,"days":10}`)
	req, _ := http.NewRequest(http.MethodPost, "/promotions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Err DB
	mockSvc.On("CreatePromo", "ERR", 10.0, 5).Return((*models.Promotion)(nil), errors.New("db err"))
	reqErr, _ := http.NewRequest(http.MethodPost, "/promotions", bytes.NewBuffer([]byte(`{"code":"ERR","discount":10.0,"days":5}`)))
	wErr := httptest.NewRecorder()
	router.ServeHTTP(wErr, reqErr)
	assert.Equal(t, http.StatusInternalServerError, wErr.Code)

	// Invalid
	reqInv, _ := http.NewRequest(http.MethodPost, "/promotions", bytes.NewBuffer([]byte(`{}`)))
	wInv := httptest.NewRecorder()
	router.ServeHTTP(wInv, reqInv)
	assert.Equal(t, http.StatusBadRequest, wInv.Code)
}
