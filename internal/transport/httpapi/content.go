package httpapi

import (
	"context"
	"fmt"
	"github.com/yohannes/kidasie-backend/internal/domain"
	"net/http"
)

type ContentReader interface {
	GetLiturgyContentBySlug(ctx context.Context, slug string) (*domain.LiturgyContent, error)
}

type contentResponse struct {
	Data *domain.LiturgyContent `json:"data"`
}

func getLiturgyContent(reader ContentReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		if slug == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{
				Error: "slug is required",
			})
			return
		}

		content, err := reader.GetLiturgyContentBySlug(r.Context(), slug)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{
				Error: "internal server error",
			})
			return
		}

		if content == nil {
			writeJSON(w, http.StatusNotFound, errorResponse{
				Error: "liturgy content not found",
			})
			return
		}

		etag := fmt.Sprintf(
			`"liturgy-%d-v%d"`,
			content.Liturgy.ID,
			content.Liturgy.ContentVersion,
		)

		w.Header().Set("ETag", etag)
		w.Header().Set("Cache-Control", "private, no-cache")

		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		writeJSON(w, http.StatusOK, contentResponse{
			Data: content,
		})
	}
}
