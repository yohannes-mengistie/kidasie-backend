package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/yohannes/kidasie-backend/internal/domain"
)

type VerseLister interface {
	ListVersesBySectionID(
		ctx context.Context,
		sectionID int64,
	) ([]domain.Verse, error)
}

type listVersesResponse struct {
	Data []domain.Verse `json:"data"`
}

func listVerses(lister VerseLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sectionID, err := strconv.ParseInt(
			r.PathValue("id"),
			10,
			64,
		)
		if err != nil || sectionID <= 0 {
			writeJSON(w, http.StatusBadRequest, errorResponse{
				Error: "section ID must be a positive integer",
			})
			return
		}

		verses, err := lister.ListVersesBySectionID(
			r.Context(),
			sectionID,
		)
		if err != nil {
			slog.Error(
				"failed to list section verses",
				"section_id",
				sectionID,
				"error",
				err,
			)

			writeJSON(w, http.StatusInternalServerError, errorResponse{
				Error: "internal server error",
			})
			return
		}

		if verses == nil {
			verses = []domain.Verse{}
		}

		writeJSON(w, http.StatusOK, listVersesResponse{
			Data: verses,
		})
	}
}
