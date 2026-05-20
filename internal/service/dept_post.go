package service

import (
	"context"
	"department/internal/models"
	"fmt"
)

func (s *DepartmentService) Add(ctx context.Context, dept models.Department) (models.Department, error) {
	newDept, err := s.repo.CreateDP(ctx, dept)
	if err != nil {
		return models.Department{}, fmt.Errorf("service add: %w", err)
	}

	return newDept, nil
}
