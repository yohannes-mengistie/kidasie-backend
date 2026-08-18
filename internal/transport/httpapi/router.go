package httpapi

import (
	"net/http"
)

type RouterDependencies struct{
	Liturgy LiturgyReader
	Section SectionLister
	Verse VerseLister
	Readiness ReadinessChecker
	Content ContentReader
}

func NewRouter(dependency RouterDependencies) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /health", health)
	if dependency.Liturgy != nil{
		router.HandleFunc("GET /api/v1/liturgies", listLiturgy(dependency.Liturgy))
		router.HandleFunc("GET /api/v1/liturgies/{slug}",getLiturgy(dependency.Liturgy))

	}

	if dependency.Section != nil{
		router.HandleFunc("GET /api/v1/liturgies/{slug}/sections",listSections(dependency.Section))
	}

	if dependency.Verse != nil{
		router.HandleFunc("GET /api/v1/sections/{id}/verses",listVerses(dependency.Verse))
	}

	if dependency.Content != nil{
		router.HandleFunc("GET /api/v1/liturgies/{slug}/content",getLiturgyContent(dependency.Content))
	}

	if dependency.Readiness != nil{
		router.HandleFunc("GET /ready", readiness(dependency.Readiness))
	}
	
	return router
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
