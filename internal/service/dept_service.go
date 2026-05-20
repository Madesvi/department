// Package service
package service

import (
	"context"
	"department/internal/models"
)

type DepartmentRepo interface {
	CreateDP(ctx context.Context, department models.Department) (models.Department, error)
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
