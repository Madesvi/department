// Package service
package service

import (
	"context"
	"department/internal/models"
)

type DepartmentRepo interface {
	GetByID(ctx context.Context, id string) (models.Employee, error)
}

type DepartmentService struct {
	repo DepartmentRepo
	// cache ...
}

func NewDepartmentService(r DepartmentRepo) *DepartmentService {
	return &DepartmentService{repo: r}
}
