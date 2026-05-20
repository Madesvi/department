package handlers

import (
	"context"
	"department/internal/models"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type DepartmentGetter interface {
	GetDepartment(ctx context.Context, id int, depth int, includeEmployee bool) (models.Department, error)
}

func GetDepartmentHandler(get DepartmentGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		depID, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Warn("invalid department id in url", "id", idStr)
			http.Error(w, "Invalid department ID", http.StatusBadRequest)
			return
		}

		query := r.URL.Query()
		depth := 1
		if dStr := query.Get("depth"); dStr != "" {
			if d, err := strconv.Atoi(dStr); err == nil {
				if d > 5 {
					d = 5
				}
				if d < 0 {
					d = 0
				}
				depth = d
			}
		}
		includeEmployees := true
		if ieStr := query.Get("include_employee"); ieStr != "" {
			includeEmployees, _ = strconv.ParseBool(ieStr)
		}

		resDep, err := get.GetDepartment(r.Context(), depID, depth, includeEmployees)
		if err != nil {
			slog.Error("get departments", "err", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := struct {
			Department models.Department
			Employee   []models.Employee
			Children   []models.Department
		}{
			Department: resDep,
			Employee:   resDep.Employee,
			Children:   resDep.Children,
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Error("error encoding response", "err", err)
		}
	}
}
