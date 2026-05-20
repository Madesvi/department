// Package service
package service

import (
	"context"
	"department/internal/models"
)

type DepartmentRepo interface {
	CreateDP(ctx context.Context, department models.Department) (models.Department, error)
}

type DepartmentService struct {
	repo DepartmentRepo
	// cache ...
}

func NewDepartmentService(r DepartmentRepo) *DepartmentService {
	return &DepartmentService{repo: r}
}

func (s *DepartmentService) Create(ctx context.Context, dept models.Department) (models.Department, error) {
	// validation... check...
	return s.repo.CreateDP(ctx, dept)
}
