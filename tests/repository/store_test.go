package repository_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
)

func TestStoreRepository_Operations(t *testing.T) {
	db := setupTestDB(t)
	// We also need to migrate Store and Staff
	db.AutoMigrate(&models.Store{}, &models.Staff{})
	
	repo := repository.NewStoreRepository(db)

	store := &models.Store{
		Name:     "Central Hub",
		Location: "Downtown",
		TaxRate:  0.1,
	}

	err := repo.Create(store)
	assert.NoError(t, err)

	found, err := repo.FindByID(store.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Central Hub", found.Name)

	staff := &models.Staff{
		StoreID: store.ID,
		Name:    "Manager Bob",
		Role:    "Manager",
	}

	err = repo.AddStaff(staff)
	assert.NoError(t, err)

	staffs, err := repo.GetStaffByStore(store.ID)
	assert.NoError(t, err)
	assert.Len(t, staffs, 1)
	assert.Equal(t, "Manager Bob", staffs[0].Name)
	
	// negative finding
	_, err = repo.FindByID(99)
	assert.Error(t, err)
}
