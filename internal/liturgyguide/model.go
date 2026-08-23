package liturgyguide

type Entry struct {
	ID       int    `json:"id"`
	Page     int    `json:"page"`
	Type     string `json:"type"`
	Language string `json:"language"`
	Text     string `json:"text"`
}

type Options struct {
	Slug           string
	Name           string
	NameAm         string
	SectionTitle   string
	SectionTitleAm string
}

type Stats struct {
	Entries        int
	MixedLanguage  int
	ReviewSegments int
}
