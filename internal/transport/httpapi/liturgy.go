package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"encoding/json"

	"github.com/yohannes/kidasie-backend/internal/domain"
)

type LiturgyLister interface{
	ListLiturgy(ctx context.Context) ([]domain.Liturgy, error)
}

type listLiturgyResponse struct{
	Data []domain.Liturgy `json:"data"`
}

type errorResponse struct{
	Error string `json:"error"`
}

func listLiturgy(lister LiturgyLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		liturgies , err := lister.ListLiturgy(r.Context())

		if err != nil{
			slog.Error("failed to list liturgies", "error", err)
			writeJson(w, http.StatusInternalServerError, errorResponse{Error: "failed to list liturgies"})
			return
		}

		if liturgies == nil{
			liturgies = []domain.Liturgy{}
		}

		writeJson(w, http.StatusOK, listLiturgyResponse{Data: liturgies})
	}
}

func writeJson(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to write JSON response", "error", err)
	}
}

