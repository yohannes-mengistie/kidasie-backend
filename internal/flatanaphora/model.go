package flatanaphora

type Entry struct {
	Page         int    `json:"page"`
	Role         string `json:"role"`
	RoleSource   string `json:"role_source"`
	EthiopicText string `json:"ethiopic_text"`
	GeezText     string `json:"geez_text"`
	AmharicText  string `json:"amharic_text"`
	EnglishText  string `json:"english_text"`
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
