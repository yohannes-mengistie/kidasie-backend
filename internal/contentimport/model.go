package contentimport

type Document struct {
	Slug     string    `json:"slug"`
	Name     string    `json:"name"`
	NameAm   string    `json:"name_am"`
	Sections []Section `json:"sections"`
}

type Section struct {
	Order   int     `json:"order"`
	Title   string  `json:"title"`
	TitleAm string  `json:"title_am"`
	Audio   *Audio  `json:"audio,omitempty"`
	Verses  []Verse `json:"verses"`
}

type Audio struct {
	URL        string `json:"url"`
	DurationMs int    `json:"duration_ms"`
	SizeBytes  int64  `json:"size_bytes"`
	MIMEType   string `json:"mime_type"`
	SHA256     string `json:"sha256"`
}

type Verse struct {
	Order    int    `json:"order"`
	TextGeez string `json:"text_geez"`
	TextAm   string `json:"text_am"`
	TextEn   string `json:"text_en,omitempty"`
	Role     string `json:"role"`
	StartMs  int    `json:"start_ms"`
	EndMs    int    `json:"end_ms"`
}
