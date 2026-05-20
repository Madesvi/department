package router

import (
	"department/internal/api/handlers"
	"net/http"
)

type RouterDeps interface {
	handlers.DeptartmentCreator
	handlers.EmployeeCreator
	handlers.DepartmentGetter
	handlers.DepartmentUpdater
}

func New(deps RouterDeps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /departments", handlers.AddDeptHandler(deps))
	mux.HandleFunc("POST /departments/{id}/employees", handlers.AddEmployeeHandler(deps))
	mux.HandleFunc("GET /departments/{id}", handlers.GetDepartmentHandler(deps))
	mux.HandleFunc("PATCH /departments/{id}", handlers.PatchDepartmentHadler(deps))
	// mux.HandleFunc("DELETE /departments/{id}", handlers.)

	return mux
}
