package repository

import (
	"github.com/user/grafana-analytics-app/internal/models"
	"gorm.io/gorm"
)

type InventoryLogRepository interface { Create(log *models.InventoryLog) error }
type invLogRepo struct{ db *gorm.DB }
func NewInventoryLogRepository(db *gorm.DB) InventoryLogRepository { return &invLogRepo{db: db} }
func (r *invLogRepo) Create(log *models.InventoryLog) error { return r.db.Create(log).Error }
