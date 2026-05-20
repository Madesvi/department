package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
)

type DepartmentDeleter interface {
	DeleteDepartment(ctx context.Context, id int, mode string, reassignID *int) error
}

func DeleteDepartmentHandler(deleter DepartmentDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		depID, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Warn("invalid id in delete url", "id", idStr)
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		query := r.URL.Query()
		mode := query.Get("mode")
		if mode == "" {
			mode = "cascade"
		}

		if mode != "cascade" && mode != "reassign" {
			http.Error(w, "Invalid mode. Use 'cascade' or 'reassign'", http.StatusBadRequest)
			return
		}

		var reassignID *int
		if mode == "reassign" {
			rIDStr := query.Get("reassign_to_department_id")
			if rIDStr == "" {
				http.Error(w, "reassign_to_department_id is required for mode=reassign", http.StatusBadRequest)
				return
			}
			val, err := strconv.Atoi(rIDStr)
			if err != nil {
				http.Error(w, "Invalid reassign_to_department_id format", http.StatusBadRequest)
				return
			}
			reassignID = &val
		}

		err = deleter.DeleteDepartment(r.Context(), depID, mode, reassignID)
		if err != nil {
			slog.Error("failed to delete department", "err", err)
			http.Error(w, "Internal server error or Target not found", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
