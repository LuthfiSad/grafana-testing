package repository_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(
		&models.Order{},
		&models.Payment{},
		&models.Shipping{},
		&models.Refund{},
		&models.OrderItem{},
		&models.Customer{},
		&models.Store{},
		&models.Staff{},
		&models.Product{},
		&models.Attribute{},
	)
	assert.NoError(t, err)
	return db
}

func TestInsertOrder(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewOrderRepository(db)

	order := &models.Order{
		CustomerID: 1,
		FinalPrice: 100.0,
		Status:     "PENDING",
	}

	err := repo.InsertOrder(order)
	assert.NoError(t, err)
	assert.NotEqual(t, uint(0), order.ID, "Order ID should be populated after insert")
}

func TestGetTotalRevenueByStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewOrderRepository(db)

	repo.InsertOrder(&models.Order{FinalPrice: 150.0, Status: "PAID"})
	repo.InsertOrder(&models.Order{FinalPrice: 200.0, Status: "PAID"})
	repo.InsertOrder(&models.Order{FinalPrice: 50.0, Status: "REFUNDED"})

	total, err := repo.GetTotalRevenueByStatus("PAID")
	assert.NoError(t, err)
	assert.Equal(t, 350.0, total, "Should calculate only PAID orders revenue")

	totalRefunded, err := repo.GetTotalRevenueByStatus("REFUNDED")
	assert.NoError(t, err)
	assert.Equal(t, 50.0, totalRefunded, "Should calculate only REFUNDED orders revenue")
	
	totalPending, err := repo.GetTotalRevenueByStatus("PENDING")
	assert.NoError(t, err)
	assert.Equal(t, 0.0, totalPending, "Should be 0 for status with no orders")
}
