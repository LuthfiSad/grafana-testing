package service

import (
	"errors"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
	"time"
)

type InventoryLogService interface { AddLog(productID uint, change int, reason string) (*models.InventoryLog, error) }
type invLogService struct{ repo repository.InventoryLogRepository }
func NewInventoryLogService(repo repository.InventoryLogRepository) InventoryLogService { return &invLogService{repo} }
func (s *invLogService) AddLog(productID uint, change int, reason string) (*models.InventoryLog, error) {
	if productID == 0 { return nil, errors.New("invalid product") }
	log := &models.InventoryLog{ProductID: productID, Change: change, Reason: reason, CreatedAt: time.Now()}
	return log, s.repo.Create(log)
}
