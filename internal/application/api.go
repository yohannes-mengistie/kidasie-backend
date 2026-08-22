package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/yohannes/kidasie-backend/internal/config"
	"github.com/yohannes/kidasie-backend/internal/database"
	"github.com/yohannes/kidasie-backend/internal/repository/postgres"
	"github.com/yohannes/kidasie-backend/internal/service"
	"github.com/yohannes/kidasie-backend/internal/transport/httpapi"
)

// RunAPI assembles and starts the Kidasie HTTP API.
func RunAPI() error {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	startupCtx, cancelStartup := context.WithTimeout(context.Background(), 10*time.Second)
	pool, err := database.OpenPostgres(startupCtx, cfg.DatabaseURL)
	cancelStartup()
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	liturgyService := service.NewLiturgyService(postgres.NewLiturgyRepo(pool))
	sectionService := service.NewSectionService(postgres.NewSectionRepo(pool))
	verseService := service.NewVerseService(postgres.NewVerseRepo(pool))
	contentService := service.NewContentService(postgres.NewContentRepo(pool))

	handler := httpapi.NewRouter(httpapi.RouterDependencies{
		Liturgy:   liturgyService,
		Section:   sectionService,
		Verse:     verseService,
		Readiness: pool,
		Content:   contentService,
	})

	server := &http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("Kidasie API is starting", "address", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve HTTP API: %w", err)
	}

	return nil
}
