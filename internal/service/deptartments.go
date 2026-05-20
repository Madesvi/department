package service

import (
	"context"
	"department/internal/models"
	"fmt"
)

func (s *DepartmentService) CreateDepartment(ctx context.Context, dept models.Department) (models.Department, error) {
	newDept, err := s.deptRepo.Create(ctx, dept)
	if err != nil {
		return models.Department{}, fmt.Errorf("service add: %w", err)
	}

	return newDept, nil
}

func (s *DepartmentService) AddEmployeeToDepartment(ctx context.Context, id int, employee models.Employee) (models.Employee, error) {
	employee.DepartmentID = id
	newEmployee, err := s.employeeRepo.CreateEmployee(ctx, employee)
	if err != nil {
		return models.Employee{}, fmt.Errorf("service add: %w", err)
	}
	return newEmployee, nil

}

func (s *DepartmentService) GetDepartment(ctx context.Context, id int, depth int, includeEmployee bool) (models.Department, error) {
	params := DepParams{
		ID:              id,
		Depth:           depth,
		IncludeEmployee: includeEmployee,
	}

	department, err := s.deptRepo.GetByID(ctx, params)
	if err != nil {
		return models.Department{}, fmt.Errorf("service get tree: %w", err)
	}
	return department, nil
}
