package httpapi

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/yohannes/kidasie-backend/internal/domain"

)

type LiturgyLister interface{
	ListLiturgies(ctx context.Context) ([]domain.Liturgy, error)
}

type listLiturgyResponse struct{
	Data []domain.Liturgy `json:"data"`
}


func listLiturgy(lister LiturgyLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		liturgies , err := lister.ListLiturgies(r.Context())

		if err != nil{
			slog.Error("failed to list liturgies", "error", err)
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to list liturgies"})
			return
		}

		if liturgies == nil{
			liturgies = []domain.Liturgy{}
		}

		writeJSON(w, http.StatusOK, listLiturgyResponse{Data: liturgies})
	}
}


