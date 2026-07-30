package main

import (
	"github.com/yohannes/kidasie-backend/internal/transport/httpapi"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	server := &http.Server{
		Addr:         ":8090",
		Handler:      httpapi.NewRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	slog.Info("Kidase Api is starting ", "address", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
