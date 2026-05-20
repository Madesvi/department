package service

import (
	"context"
	"department/internal/models"
	"fmt"
)

func (s *DepartmentService) CreateDepartment(ctx context.Context, dept models.Department) (models.Department, error) {
	newDept, err := s.deptRepo.CreateDP(ctx, dept)
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
