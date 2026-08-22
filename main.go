package main

import (
	"log/slog"
	"os"

	"github.com/yohannes/kidasie-backend/internal/application"
)

func main() {
	if err := application.RunAPI(); err != nil {
		slog.Error("Kidasie API stopped", "error", err)
		os.Exit(1)
	}
}
