package main

import (
	"department/internal/api/router"
	"department/internal/repositories/postgresql"
	"department/internal/service"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env files", "err", err)
	}

	// Load logger
	logLevel := os.Getenv("LOG_LEVEL")

	var level slog.Level
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: level}

	if os.Getenv("APP_ENV") == "development" {
		handler = slog.NewTextHandler(os.Stderr, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	// Load logger

	slog.Debug("log level set", "value", logLevel)

	db, err := postgresql.ConnectDB()
	if err != nil {
		slog.Error("database connection failed", "err", err)
		os.Exit(1)
	}
	deptRepo := postgresql.NewDepartmentRepo(db)
	empRepo := postgresql.NewEmployeeRepo(db)
	srv := service.NewDepartmentService(deptRepo, empRepo)

	r := router.New(srv)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	// Create custom server
	server := &http.Server{
		Addr:    port,
		Handler: r,
	}

	slog.Info("Server is running on port", "port", port)
	err = server.ListenAndServe()
	if err != nil {
		slog.Error("Error starting the server", "err", err)
	}
}
