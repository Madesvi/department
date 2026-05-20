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

func (r DepartmentProvider) DeleteCascade(ctx context.Context, id int) error {
	err := r.db.WithContext(ctx).Delete(&models.Department{}, id).Error
	if err != nil {
		return fmt.Errorf("db delete cascade: %w", err)
	}
	return nil
}

func (r DepartmentProvider) DeleteWithReassign(ctx context.Context, id int, newDeptID int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var check models.Department
		if err := tx.First(&check, newDeptID).Error; err != nil {
			return fmt.Errorf("target department %d not found: %w", newDeptID, err)
		}

		err := tx.Model(&models.Employee{}).
			Where("department_id = ?", id).
			Update("department_id", newDeptID).Error
		if err != nil {
			return fmt.Errorf("failed to reassign employees: %w", err)
		}

		err = tx.Model(&models.Department{}).
			Where("parent_id = ?", id).
			Update("parent_id", newDeptID).Error
		if err != nil {
			return fmt.Errorf("failed to reassign sub-departments: %w", err)
		}

		if err := tx.Delete(&models.Department{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete department: %w", err)
		}

		return nil
	})
}

func (r DepartmentProvider) IsNameDuplicate(ctx context.Context, parentID *int, name string, excludeID int) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Department{}).Where("name = ?", name)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
