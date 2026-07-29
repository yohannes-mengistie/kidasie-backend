package memory

import (
	"github.com/yohannes/kidasie-backend/internal/domain"
)

type LiturgyRepo struct {
	liturgies []domain.Liturgy
	sections  map[int64]*domain.Section
}

type SeedFile struct {
	Liturgies []domain.Liturgy `json:"liturgies"`
	Sections  []domain.Section `json:"sections"`
}
