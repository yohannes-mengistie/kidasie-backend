package flatanaphora

type Entry struct {
	Page                 int     `json:"page"`
	SourcePage           int     `json:"source_page,omitempty"`
	Kind                 string  `json:"kind,omitempty"`
	Role                 string  `json:"role"`
	RoleSource           string  `json:"role_source,omitempty"`
	EthiopicText         string  `json:"ethiopic_text,omitempty"`
	GeezText             string  `json:"geez_text,omitempty"`
	TextGeez             string  `json:"text_geez,omitempty"`
	AmharicText          string  `json:"amharic_text,omitempty"`
	EnglishText          string  `json:"english_text,omitempty"`
	OriginalEthiopicText string  `json:"original_ethiopic_text,omitempty"`
	SeparationConfidence float64 `json:"separation_confidence,omitempty"`
	SeparationNote       string  `json:"separation_note,omitempty"`
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
