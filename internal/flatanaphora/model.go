package flatanaphora

type Entry struct {
	Page                 int     `json:"page"`
	Number               int     `json:"number,omitempty"`
	SourcePage           int     `json:"source_page,omitempty"`
	Kind                 string  `json:"kind,omitempty"`
	Role                 string  `json:"role"`
	RoleSource           string  `json:"role_source,omitempty"`
	EthiopicText         string  `json:"ethiopic_text,omitempty"`
	GeezText             string  `json:"geez_text,omitempty"`
	TextGeez             string  `json:"text_geez,omitempty"`
	AmharicText          string  `json:"amharic_text,omitempty"`
	TextAmharic          string  `json:"text_amharic,omitempty"`
	EnglishText          string  `json:"english_text,omitempty"`
	TextEnglish          string  `json:"text_english,omitempty"`
	OriginalEthiopicText string  `json:"original_ethiopic_text,omitempty"`
	SeparationConfidence float64 `json:"separation_confidence,omitempty"`
	SeparationNote       string  `json:"separation_note,omitempty"`
	Parts                []Part  `json:"parts,omitempty"`
	NeedsReview          *bool   `json:"needs_review,omitempty"`
}

type Part struct {
	Role        string `json:"role"`
	TextGeez    string `json:"text_geez,omitempty"`
	TextAmharic string `json:"text_amharic,omitempty"`
	TextEnglish string `json:"text_english,omitempty"`
}

type Options struct {
	Slug           string
	Name           string
	NameAm         string
	SectionTitle   string
	SectionTitleAm string
}

type Stats struct {
	Entries         int
	AmbiguousText   int
	ReviewSegments  int
	UntimedSegments int
}
