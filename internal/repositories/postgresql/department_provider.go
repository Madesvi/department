package postgresql

import (
	"context"
	"department/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type DepartmentProvider struct {
	db *gorm.DB
}

func NewDepartmentRepo(db *gorm.DB) *DepartmentProvider {
	return &DepartmentProvider{db: db}
}

func (r DepartmentProvider) CreateDP(ctx context.Context, department models.Department) (models.Department, error) {
	if err := r.db.WithContext(ctx).Create(&department).Error; err != nil {
		return models.Department{}, fmt.Errorf("create department: %w", err)
	}
	return department, nil
}
