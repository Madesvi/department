package handlers

import (
	"context"
	"department/internal/models"
	"encoding/json"
	"log/slog"
	"net/http"
)

type DeptartmentCreator interface {
	CreateDepartment(ctx context.Context, dept models.Department) (models.Department, error)
}

type CreateRequest struct {
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id"`
}

func AddDeptHandler(create DeptartmentCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRequest

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Warn("failed to decode body", "err", err)
			http.Error(w, "Invalid request Body", http.StatusBadRequest)
			return
		}

		if req.Name == "" {
			http.Error(w, "Faild 'name' is required", http.StatusBadRequest)
			return
		}

		srvDept := models.Department{
			Name:     req.Name,
			ParentID: req.ParentID,
		}

		newDept, err := create.CreateDepartment(r.Context(), srvDept)
		if err != nil {
			slog.Error("add departments", "err", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := struct {
			Status string `json:"status"`
			Data   models.Department
		}{
			Status: "success",
			Data:   newDept,
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Error("error encoding response", "err", err)
		}
	}
}
