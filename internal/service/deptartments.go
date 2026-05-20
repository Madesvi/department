package service

import (
	"context"
	"department/internal/models"
	"errors"
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

func (s *DepartmentService) UpdateDepartment(ctx context.Context, id int, updates map[string]any) (models.Department, error) {
	if parentID, ok := updates["parent_id"]; ok && parentID != nil {
		if pIDFloat, ok := parentID.(float64); ok {
			if int(pIDFloat) == id {
				return models.Department{}, errors.New("a department cannot be its own parent")
			}
		}
	}
	dept, err := s.deptRepo.Update(ctx, id, updates)
	if err != nil {
		return models.Department{}, fmt.Errorf("service update: %w", err)
	}
	return dept, nil
}

func (s *DepartmentService) DeleteDepartment(ctx context.Context, id int, mode string, reassignID *int) error {
	if mode == "reassign" {
		if reassignID == nil {
			return fmt.Errorf("reassign_to_department_id is required when mode=reassign")
		}
		if *reassignID == id {
			return fmt.Errorf("cannot reassign employees to the department being deleted")
		}
		return s.deptRepo.DeleteWithReassign(ctx, id, *reassignID)
	}

	return s.deptRepo.DeleteCascade(ctx, id)
}
