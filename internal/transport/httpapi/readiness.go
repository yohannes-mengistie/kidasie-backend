package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type ReadinessChecker interface {
	Ping(ctx context.Context) error
}

type readinessResponse struct {
	Status string `json:"status"`
}

func readiness(checker ReadinessChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(
			r.Context(),
			2*time.Second,
		)
		defer cancel()

		if err := checker.Ping(ctx); err != nil {
			slog.Error(
				"readiness check failed",
				"error",
				err,
			)

			writeJSON(w, http.StatusServiceUnavailable, readinessResponse{
				Status: "not_ready",
			})
			return
		}

		writeJSON(w, http.StatusOK, readinessResponse{
			Status: "ready",
		})
	}
}
