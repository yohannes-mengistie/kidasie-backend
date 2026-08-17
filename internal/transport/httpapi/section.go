package httpapi

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/yohannes/kidasie-backend/internal/domain"
)


type SectionLister interface{
	ListSectionsByLiturgySlug(ctx context.Context , slug string)([]domain.Section , error)
}

type listSectionReponse struct{
	Data []domain.Section `json:"data"`
}

func listSections(lister SectionLister) http.HandlerFunc{
	return func(w http.ResponseWriter , r *http.Request){
		slug := r.PathValue("slug")

		if slug == ""{
			writeJson(w, http.StatusBadRequest,errorResponse{
				Error : "liturgy slug is required",
			})

			return
		}
		sections,err := lister.ListSectionsByLiturgySlug(r.Context() , slug,)

		if err != nil {
			slog.Error(
				"failed to list liturgy section",
				"error",
				err,
			)
			writeJson(w , http.StatusInternalServerError,errorResponse{
				Error: "internal server error",
			})
			return
		}

		if sections == nil {
			sections = []domain.Section{}
		}

		writeJson(w,http.StatusOK,listSectionReponse{
			Data:sections,
		})


	}
}