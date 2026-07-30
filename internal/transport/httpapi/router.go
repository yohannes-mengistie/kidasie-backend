package httpapi

import (
	"net/http"
)

func NewRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /health", health)

	return router
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
