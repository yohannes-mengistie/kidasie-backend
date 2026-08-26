package domain

import "testing"

func TestAnnouncementInputNormalizeAndValidate(t *testing.T) {
	t.Parallel()

	input := AnnouncementInput{
		Slug:        " new-content ",
		TitleAm:     " አዲስ ይዘት ",
		TitleEn:     " New content ",
		BodyAm:      " አዲስ ይዘት ተጨምሯል። ",
		BodyEn:      " New content is available. ",
		Kind:        " content ",
		ActionType:  " open_liturgy ",
		ActionValue: " anaphora-of-st-mary ",
		Status:      AnnouncementStatusPublished,
	}

	if err := input.NormalizeAndValidate(); err != nil {
		t.Fatalf("NormalizeAndValidate() error = %v", err)
	}
	if input.Slug != "new-content" {
		t.Fatalf("Slug = %q, want %q", input.Slug, "new-content")
	}
	if input.ActionValue != "anaphora-of-st-mary" {
		t.Fatalf("ActionValue = %q", input.ActionValue)
	}
}

func TestAnnouncementInputRejectsUnsafeDownloadURL(t *testing.T) {
	t.Parallel()

	input := validAnnouncementInput()
	input.ActionType = "download_apk"
	input.ActionValue = "http://example.com/app.apk"

	if err := input.NormalizeAndValidate(); err == nil {
		t.Fatal("NormalizeAndValidate() error = nil, want HTTPS validation error")
	}
}

func TestAnnouncementInputRequiresActionValue(t *testing.T) {
	t.Parallel()

	input := validAnnouncementInput()
	input.ActionType = "open_liturgy"
	input.ActionValue = ""

	if err := input.NormalizeAndValidate(); err == nil {
		t.Fatal("NormalizeAndValidate() error = nil, want missing value error")
	}
}

func validAnnouncementInput() AnnouncementInput {
	return AnnouncementInput{
		Slug:       "announcement",
		TitleAm:    "ርዕስ",
		TitleEn:    "Title",
		BodyAm:     "መልእክት",
		BodyEn:     "Message",
		Kind:       "general",
		ActionType: "none",
		Status:     AnnouncementStatusDraft,
	}
}
