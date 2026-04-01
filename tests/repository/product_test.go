package repository_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
)

func TestProductRepository_Operations(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&models.Product{}, &models.Attribute{})
	
	repo := repository.NewProductRepository(db)

	prod := &models.Product{
		Name:  "Test Item",
		Price: 100.0,
		Stock: 50,
		Attributes: []models.Attribute{
			{Key: "Color", Value: "Red"},
		},
	}

	err := repo.Create(prod)
	assert.NoError(t, err)
	assert.NotEqual(t, uint(0), prod.ID)

	found, err := repo.FindByID(prod.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Test Item", found.Name)
	assert.Len(t, found.Attributes, 1)

	// Update stock
	err = repo.UpdateStock(prod.ID, -10)
	assert.NoError(t, err)
	
	updated, _ := repo.FindByID(prod.ID)
	assert.Equal(t, 40, updated.Stock)
	
	// negative finding
	_, err = repo.FindByID(99)
	assert.Error(t, err)
}
