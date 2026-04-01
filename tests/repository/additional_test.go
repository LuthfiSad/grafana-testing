package repository_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
)

func TestPromotionRepository(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&models.Promotion{})
	repo := repository.NewPromotionRepository(db)

	err := repo.Create(&models.Promotion{Code: "PROMO10", Discount: 10.0})
	assert.NoError(t, err)

	promo, err := repo.FindByCode("PROMO10")
	assert.NoError(t, err)
	assert.Equal(t, 10.0, promo.Discount)

	_, err = repo.FindByCode("INVALID")
	assert.Error(t, err)
}

func TestPaymentRepository(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&models.Payment{})
	repo := repository.NewPaymentRepository(db)

	err := repo.Create(&models.Payment{OrderID: 1, Method: "CC"})
	assert.NoError(t, err)

	pay, err := repo.FindByOrderID(1)
	assert.NoError(t, err)
	assert.Equal(t, "CC", pay.Method)

	_, err = repo.FindByOrderID(99)
	assert.Error(t, err)
}

func TestShippingRepository(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&models.Shipping{})
	repo := repository.NewShippingRepository(db)
	err := repo.Create(&models.Shipping{OrderID: 1, Carrier: "FedEx"})
	assert.NoError(t, err)
}

func TestReviewRepository(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&models.Review{})
	repo := repository.NewReviewRepository(db)
	err := repo.Create(&models.Review{ProductID: 1, Rating: 5})
	assert.NoError(t, err)
}

func TestRefundRepository(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&models.Refund{})
	repo := repository.NewRefundRepository(db)
	err := repo.Create(&models.Refund{OrderID: 1, Amount: 100.0})
	assert.NoError(t, err)
}

func TestInventoryLogRepository(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&models.InventoryLog{})
	repo := repository.NewInventoryLogRepository(db)
	err := repo.Create(&models.InventoryLog{ProductID: 1, Change: 10})
	assert.NoError(t, err)
}

func TestAttributeRepository(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&models.Attribute{})
	repo := repository.NewAttributeRepository(db)
	err := repo.Create(&models.Attribute{ProductID: 1, Key: "Size", Value: "M"})
	assert.NoError(t, err)
}

func TestOrderItemRepository(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&models.OrderItem{})
	repo := repository.NewOrderItemRepository(db)
	err := repo.Create(&models.OrderItem{OrderID: 1, ProductID: 1, Quantity: 2})
	assert.NoError(t, err)
}
