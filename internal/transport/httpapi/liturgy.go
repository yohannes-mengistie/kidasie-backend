package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"errors"

	"github.com/yohannes/kidasie-backend/internal/domain"

)

type LiturgyReader interface{
	ListLiturgies(ctx context.Context) ([]domain.Liturgy, error)
	GetLiturgyBySlug(
  		ctx context.Context,
  		slug string,
  	) (*domain.Liturgy, error)

}

type listLiturgyResponse struct{
	Data []domain.Liturgy `json:"data"`
}

type getLiturgyResponse struct{
	Data domain.Liturgy `json:"data"`
}


func listLiturgy(lister LiturgyReader) http.HandlerFunc {
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

func getLiturgy(reader LiturgyReader) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
  		slug := r.PathValue("slug")
  		if slug == "" {
  			writeJSON(w, http.StatusBadRequest, errorResponse{
  				Error: "liturgy slug is required",
  			})
  			return
  		}

  		liturgy, err := reader.GetLiturgyBySlug(
  			r.Context(),
  			slug,
  		)
  		if errors.Is(err, domain.ErrNotFound) {
  			writeJSON(w, http.StatusNotFound, errorResponse{
  				Error: "liturgy not found",
  			})
  			return
  		}
  		if err != nil {
  			slog.Error(
  				"failed to get liturgy",
  				"slug",
  				slug,
  				"error",
  				err,
  			)

  			writeJSON(w, http.StatusInternalServerError, errorResponse{
  				Error: "internal server error",
  			})
  			return
  		}

  		writeJSON(w, http.StatusOK, getLiturgyResponse{
  			Data: *liturgy,
  		})
  	}

}


