package service

import (
	"context"
	"department/internal/models"
	"errors"
	"fmt"
	"strings"
)

func validateString(val string, min, max int) (string, error) {
	trimmed := strings.TrimSpace(val)
	length := len([]rune(trimmed))
	if length < min || length > max {
		return "", fmt.Errorf("Length should be from %d to %d", min, max)
	}
	return trimmed, nil
}

func (s *DepartmentService) CreateDepartment(ctx context.Context, dept models.Department) (models.Department, error) {

	cleanedName, err := validateString(dept.Name, 1, 200)
	if err != nil {
		return models.Department{}, fmt.Errorf("clean name: %w", err)
	}
	dept.Name = cleanedName

	isDuplicate, err := s.deptRepo.IsNameDuplicate(ctx, dept.ParentID, dept.Name, 0)
	if err != nil {
		return models.Department{}, fmt.Errorf("check unique name failed: %w", err)
	}
	if isDuplicate {
		return models.Department{}, errors.New("409, name must be unique")
	}

	newDept, err := s.deptRepo.Create(ctx, dept)
	if err != nil {
		return models.Department{}, fmt.Errorf("service add: %w", err)
	}

	return newDept, nil
}

func (s *DepartmentService) AddEmployeeToDepartment(ctx context.Context, id int, employee models.Employee) (models.Employee, error) {
	// Check ID department
	_, err := s.deptRepo.GetByID(ctx, DepParams{ID: id})
	if err != nil {
		return models.Employee{}, errors.New("404: cannot create employee")
	}

	cleanedName, err := validateString(employee.FullName, 1, 200)
	if err != nil {
		return models.Employee{}, fmt.Errorf("clean full_name: %w", err)
	}
	employee.FullName = cleanedName

	cleanPosition, err := validateString(employee.Position, 1, 200)
	if err != nil {
		return models.Employee{}, fmt.Errorf("position: %w", err)
	}
	employee.Position = cleanPosition

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
	// Check if name was updated
	if nameRaw, ok := updates["name"]; ok {
		if nameStr, ok := nameRaw.(string); ok {
			cleanedName, err := validateString(nameStr, 1, 200)
			if err != nil {
				return models.Department{}, fmt.Errorf("name: %w", err)
			}
			updates["name"] = cleanedName
		}
	}

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
