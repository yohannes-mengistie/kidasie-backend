package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yohannes/kidasie-backend/internal/domain"
)

type announcementReaderStub struct {
	announcements []domain.Announcement
	err           error
}

func (stub announcementReaderStub) ListPublishedAnnouncements(
	context.Context,
) ([]domain.Announcement, error) {
	return stub.announcements, stub.err
}

func TestListAnnouncementsReturnsFeedAndETag(t *testing.T) {
	t.Parallel()

	publishedAt := time.Date(2026, time.August, 26, 9, 0, 0, 0, time.UTC)
	router := NewRouter(RouterDependencies{
		Announcement: announcementReaderStub{announcements: []domain.Announcement{
			{
				ID:          1,
				Slug:        "new-content",
				Version:     2,
				TitleAm:     "አዲስ ይዘት",
				TitleEn:     "New content",
				BodyAm:      "ይዘት ተጨምሯል።",
				BodyEn:      "Content is available.",
				Kind:        "content",
				Action:      domain.AnnouncementAction{Type: "none"},
				PublishedAt: &publishedAt,
				UpdatedAt:   publishedAt,
			},
		}},
	})

	request := httptest.NewRequest(http.MethodGet, "/api/v1/announcements", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if recorder.Header().Get("ETag") == "" {
		t.Fatal("ETag is empty")
	}
	if !strings.Contains(recorder.Body.String(), `"slug":"new-content"`) {
		t.Fatalf("body = %s", recorder.Body.String())
	}

	request = httptest.NewRequest(http.MethodGet, "/api/v1/announcements", nil)
	request.Header.Set("If-None-Match", recorder.Header().Get("ETag"))
	notModified := httptest.NewRecorder()
	router.ServeHTTP(notModified, request)
	if notModified.Code != http.StatusNotModified {
		t.Fatalf("conditional status = %d, want %d", notModified.Code, http.StatusNotModified)
	}
}
