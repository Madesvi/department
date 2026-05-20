package router

import (
	"department/internal/api/handlers"
	"net/http"
)

type RouterDeps interface {
	handlers.DeptProvider
}

func New(deps RouterDeps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /departments", handlers.AddDeptHandler(deps))
	// mux.HandleFunc("POST /departments/{id}/employees", handlers.)
	// mux.HandleFunc("GET /departments/{id}", handlers.)
	// mux.HandleFunc("PATCH /departments/{id}", handlers.)
	// mux.HandleFunc("DELETE /departments/{id}", handlers.)

	return mux
}
