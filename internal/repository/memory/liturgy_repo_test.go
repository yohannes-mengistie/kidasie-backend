package memory

import (
	"context"
	"github.com/yohannes/kidasie-backend/internal/domain"
	"testing"
)

func TestLiturgyRepo_ListLiturgies(t *testing.T) {
	input := []domain.Liturgy{
		{
			ID:     1,
			Name:   "Anaphora of the Apostles",
			NameAm: "የሐዋርያት ቅዳሴ",
		},
		{
			ID:     2,
			Name:   "Anaphora of Our Lord",
			NameAm: "የጌታችን ቅዳሴ",
		},
	}

	repo := NewLiturgyRepo(input, nil)
	got, err := repo.ListLiturgies(context.Background())
	if err != nil {
		t.Fatalf("ListLiturgies() error = %v", err)
	}
	if len(got) != len(input) {
		t.Fatalf("expected %d liturgies, got %d", len(input), len(got))
	}

	if got[0].ID != input[0].ID {
		t.Errorf("expected first ID %d, got %d", input[0].ID, got[0].ID)
	}

	got[0].Name = "Changed by caller"

	again, err := repo.ListLiturgies(context.Background())
	if err != nil {
		t.Fatalf("second ListLiturgies() returned error: %v", err)
	}

	if again[0].Name == "Changed by caller" {
		t.Error("repository data was modified by the caller")
	}

}
