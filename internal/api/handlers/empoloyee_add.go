package handlers

import (
	"context"
	"department/internal/models"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type EmployeeCreator interface {
	AddEmployeeToDepartment(ctx context.Context, id int, employee models.Employee) (models.Employee, error)
}

type CreateEmployeeRequest struct {
	FullName string     `gorm:"not null" json:"full_name"`
	Position string     `gorm:"not null" json:"position"`
	HiredAt  *time.Time `json:"hired_at,omitempty"`
}

func AddEmployeeHandler(create EmployeeCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := r.PathValue("id")
		depID, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Warn("invalid department id in url", "id", idStr)
			http.Error(w, "Invalid department ID", http.StatusBadRequest)
			return
		}

		var req CreateEmployeeRequest
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Warn("failed to decode body", "err", err)
			http.Error(w, "Invalid request Body", http.StatusBadRequest)
			return
		}

		if req.FullName == "" || req.Position == "" {
			http.Error(w, "Faild 'full name' and 'position' is required", http.StatusBadRequest)
			return
		}

		srvEmp := models.Employee{
			FullName: req.FullName,
			Position: req.Position,
			HiredAt:  req.HiredAt,
		}

		newEmployee, err := create.AddEmployeeToDepartment(r.Context(), depID, srvEmp)
		if err != nil {
			slog.Error("add departments", "err", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := struct {
			Status string `json:"status"`
			Data   models.Employee
		}{
			Status: "success",
			Data:   newEmployee,
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Warn("error encoding response", "err", err)
		}
	}
}
