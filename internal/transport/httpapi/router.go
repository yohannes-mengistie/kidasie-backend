package httpapi

import (
	"net/http"
)

type RouterDependencies struct{
	Liturgy LiturgyLister
	Section SectionLister
}

func NewRouter(dependency RouterDependencies) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /health", health)
	if dependency.Liturgy != nil{
		router.HandleFunc("GET /api/v1/liturgies", listLiturgy(dependency.Liturgy))
	}

	if dependency.Section != nil{
		router.HandleFunc("GET /api/v1/liturgies/{slug}/sections",listSections(dependency.Section))
	}
	
	return router
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
