package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthReturnsOk(t *testing.T) {
	router := NewRouter(RouterDependencies{})
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if got := recorder.Header().Get("Content-Type"); got !=
		"application/json" {
		t.Errorf("expected JSON content type, got %q", got)
	}

	const expectedBody = `{"status":"ok"}`
	if recorder.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody,
			recorder.Body.String())
	}

}
