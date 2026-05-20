package postgresql

import (
	"context"
	"department/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type EmployeeProvider struct {
	db *gorm.DB
}

func NewEmployeeRepo(db *gorm.DB) *EmployeeProvider {
	return &EmployeeProvider{db: db}
}

func (r EmployeeProvider) CreateEmployee(ctx context.Context, employee models.Employee) (models.Employee, error) {
	if err := r.db.WithContext(ctx).Create(&employee).Error; err != nil {
		return models.Employee{}, fmt.Errorf("create employee: %w", err)
	}
	return employee, nil
}
