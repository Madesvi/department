// Package service
package service

import (
	"context"
	"department/internal/models"
)

type DepParams struct {
	ID              int
	Depth           int
	IncludeEmployee bool
}

type DepartmentRepo interface {
	Create(ctx context.Context, department models.Department) (models.Department, error)
	GetByID(ctx context.Context, params DepParams) (models.Department, error)
	Update(ctx context.Context, id int, updates map[string]any) (models.Department, error)
	DeleteCascade(ctx context.Context, id int) error
	DeleteWithReassign(ctx context.Context, id int, newDeptID int) error
	IsNameDuplicate(ctx context.Context, parentID *int, name string, excludeID int) (bool, error)
}

type EmployeeRepo interface {
	CreateEmployee(ctx context.Context, employee models.Employee) (models.Employee, error)
}

type DepartmentService struct {
	deptRepo     DepartmentRepo
	employeeRepo EmployeeRepo
}

func NewDepartmentService(dr DepartmentRepo, er EmployeeRepo) *DepartmentService {
	return &DepartmentService{deptRepo: dr, employeeRepo: er}
}
