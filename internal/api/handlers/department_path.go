package handlers

import (
	"context"
	"department/internal/models"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type DepartmentUpdater interface {
	UpdateDepartment(ctx context.Context, id int, updates map[string]any) (models.Department, error)
}

func PatchDepartmentHadler(update DepartmentUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		depID, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Warn("invalid id in patch url", "id", idStr)
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var updates map[string]any
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			slog.Warn("failed to decode patch body", "err", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		resDep, err := update.UpdateDepartment(r.Context(), depID, updates)
		if err != nil {
			slog.Error("failed to update department", "err", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(resDep)
		if err != nil {
			slog.Warn("error encoding response", "err", err)
		}
	}
}
