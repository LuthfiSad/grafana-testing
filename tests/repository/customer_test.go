package repository_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
)

func TestCustomerRepository_CreateAndFind(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewCustomerRepository(db)

	cust := &models.Customer{
		Name:    "John Doe",
		Email:   "john@test.com",
		Segment: "VIP",
	}

	err := repo.Create(cust)
	assert.NoError(t, err)
	assert.NotEqual(t, uint(0), cust.ID)

	found, err := repo.FindByID(cust.ID)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", found.Name)

	_, err = repo.FindByID(9999)
	assert.Error(t, err)
}

func TestCustomerRepository_UpdateLoyalty(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewCustomerRepository(db)

	cust := &models.Customer{Name: "Jane", Email: "jane@test.com", LoyaltyPoints: 10}
	repo.Create(cust)

	err := repo.UpdateLoyalty(cust.ID, 15)
	assert.NoError(t, err)

	updated, _ := repo.FindByID(cust.ID)
	assert.Equal(t, 25, updated.LoyaltyPoints)
}
