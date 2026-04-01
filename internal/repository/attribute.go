package repository

import (
	"github.com/user/grafana-analytics-app/internal/models"
	"gorm.io/gorm"
)

type AttributeRepository interface { Create(attr *models.Attribute) error }
type attrRepo struct{ db *gorm.DB }
func NewAttributeRepository(db *gorm.DB) AttributeRepository { return &attrRepo{db: db} }
func (r *attrRepo) Create(attr *models.Attribute) error { return r.db.Create(attr).Error }
