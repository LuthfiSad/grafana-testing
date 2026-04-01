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

type MockPaymentSvc struct{ mock.Mock }
func (m *MockPaymentSvc) ProcessPayment(oid uint, meth string) (*models.Payment, error) {
	args := m.Called(oid, meth)
	if args.Get(0) != nil { return args.Get(0).(*models.Payment), args.Error(1) }
	return nil, args.Error(1)
}

func TestPaymentHandler_ProcessPayment(t *testing.T) {
	mockSvc := new(MockPaymentSvc)
	h := handler.NewPaymentHandler(mockSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/payments", h.HandleProcessPayment)

	mockSvc.On("ProcessPayment", uint(1), "CC").Return(&models.Payment{Method: "CC"}, nil)
	body := []byte(`{"order_id":1,"method":"CC"}`)
	req, _ := http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Err
	mockSvc.On("ProcessPayment", uint(2), "CC").Return((*models.Payment)(nil), errors.New("db err"))
	reqErr, _ := http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer([]byte(`{"order_id":2,"method":"CC"}`)))
	wErr := httptest.NewRecorder()
	router.ServeHTTP(wErr, reqErr)
	assert.Equal(t, http.StatusInternalServerError, wErr.Code)

	// Inv
	reqInv, _ := http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer([]byte(`{}`)))
	wInv := httptest.NewRecorder()
	router.ServeHTTP(wInv, reqInv)
	assert.Equal(t, http.StatusBadRequest, wInv.Code)
}
