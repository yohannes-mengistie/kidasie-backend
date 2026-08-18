package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/yohannes/kidasie-backend/internal/config"
	"github.com/yohannes/kidasie-backend/internal/database"
	"github.com/yohannes/kidasie-backend/internal/repository/postgres"
	"github.com/yohannes/kidasie-backend/internal/service"
	"github.com/yohannes/kidasie-backend/internal/transport/httpapi"
)


func main(){
	if err := run(); err != nil{
		slog.Error("Kidasie API is stoped","error", err)
		os.Exit(1)
	}
}

func run() error{
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		slog.Error("Invalid configuration", "error", err)
		os.Exit(1)
	}
	startupCtx , cancelStartUp := context.WithTimeout(context.Background(),10*time.Second,)

	pool,err := database.OpenPostgres(startupCtx,cfg.DatabaseURL)
	cancelStartUp()
	if err != nil{
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	liturgyRepository := postgres.NewLiturgyRepo(pool)
	sectionRepository := postgres.NewSectionRepo(pool)
	verseRepository := postgres.NewVerseRepo(pool)
	contentRepository := postgres.NewContentRepo(pool)
	liturgyService := service.NewLiturgyService(liturgyRepository)
	sectionService := service.NewSectionService(sectionRepository)
	verseService := service.NewVerseService(verseRepository)
	contentService := service.NewContentService(contentRepository)
	
	handler := httpapi.NewRouter(httpapi.RouterDependencies{Liturgy: liturgyService , Section: sectionService , Verse :verseService , Readiness: pool, Content: contentService})

	server := &http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	slog.Info("Kidase Api is starting ", "address", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}

	return nil
}
