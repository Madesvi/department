package postgresql

import (
	"context"
	"department/internal/models"
	"department/internal/service"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DepartmentProvider struct {
	db *gorm.DB
}

func NewDepartmentRepo(db *gorm.DB) *DepartmentProvider {
	return &DepartmentProvider{db: db}
}

func (r DepartmentProvider) Create(ctx context.Context, department models.Department) (models.Department, error) {
	if err := r.db.WithContext(ctx).Create(&department).Error; err != nil {
		return models.Department{}, fmt.Errorf("create department: %w", err)
	}
	return department, nil
}

func (r DepartmentProvider) GetByID(ctx context.Context, params service.DepParams) (models.Department, error) {
	var dept models.Department
	query := r.db.WithContext(ctx)

	if params.Depth > 0 {
		level := make([]string, params.Depth)
		for i := 0; i < params.Depth; i++ {
			level[i] = "Children"
		}
		preloadPath := strings.Join(level, ".")
		query = query.Preload(preloadPath)
	}

	if params.IncludeEmployee {
		query = query.Preload("Employee", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		})
	}
	if err := query.First(&dept, params.ID).Error; err != nil {
		return models.Department{}, err
	}
	return dept, nil

}

func (r DepartmentProvider) Update(ctx context.Context, id int, updates map[string]any) (models.Department, error) {
	var updatedDept models.Department

	err := r.db.WithContext(ctx).
		Model(&updatedDept).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updates).Error

	if err != nil {
		return models.Department{}, fmt.Errorf("db update department: %w", err)
	}

	return updatedDept, nil
}
