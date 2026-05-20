package main

import (
	postgre "department/internal/repositories/postgresql"
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

	db, err := postgre.ConnectDB()
	if err != nil {
		slog.Error("database connection failed", "err", err)
		os.Exit(1)
	}

	port := os.Getenv("SERVER_PORT")

	secureMux := utils.ApplyMiddleWares(router, mw.SecurityHeaders)

	// Create custom server
	server := &http.Server{
		Addr:    port,
		Handler: secureMux,
		// TLSConfig: tlsConfig,
	}

	slog.Info("Server is running on port", "port", port)
	err = server.ListenAndServe()
	if err != nil {
		slog.Error("Error starting the server", "err", err)
	}
}
