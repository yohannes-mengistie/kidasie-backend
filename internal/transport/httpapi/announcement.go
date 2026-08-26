package httpapi

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/yohannes/kidasie-backend/internal/domain"
)

type AnnouncementReader interface {
	ListPublishedAnnouncements(ctx context.Context) ([]domain.Announcement, error)
}

type listAnnouncementsResponse struct {
	Data []domain.Announcement `json:"data"`
}

func listAnnouncements(reader AnnouncementReader) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		announcements, err := reader.ListPublishedAnnouncements(request.Context())
		if err != nil {
			slog.Error("failed to list announcements", "error", err)
			writeJSON(w, http.StatusInternalServerError, errorResponse{
				Error: "failed to list announcements",
			})
			return
		}

		if announcements == nil {
			announcements = []domain.Announcement{}
		}

		hasher := sha256.New()
		for _, announcement := range announcements {
			_, _ = fmt.Fprintf(
				hasher,
				"%d:%d;",
				announcement.ID,
				announcement.Version,
			)
		}
		etag := fmt.Sprintf(`"announcements-%x"`, hasher.Sum(nil)[:12])

		w.Header().Set("ETag", etag)
		w.Header().Set("Cache-Control", "private, no-cache")
		if request.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		writeJSON(w, http.StatusOK, listAnnouncementsResponse{
			Data: announcements,
		})
	}
}
