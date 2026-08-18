package domain

type Verse struct {
	ID       int64  `json:"id"`
	Order    int    `json:"order"`
	TextGeez string `json:"text_geez"`
	TextAm   string `json:"text_am"`
	TextEn   string `json:"text_en"`
	Role     string `json:"role"`
	StartMs  int    `json:"start_ms"`
	EndMs    int    `json:"end_ms"`
}

type Section struct {
	ID              int64   `json:"id"`
	LiturgyID       int64   `json:"liturgy_id"`
	Order           int     `json:"order"`
	Title           string  `json:"title"`
	TitleAm         string  `json:"title_am"`
	AudioURL        string  `json:"audio_url,omitempty"`
	AudioDurationMs int     `json:"audio_duration_ms,omitempty"`
	AudioSizeBytes  int64   `json:"audio_size_bytes,omitempty"`
	AudioMIMEType   string  `json:"audio_mime_type,omitempty"`
	AudioSHA256     string  `json:"audio_sha256,omitempty"`
	Verses          []Verse `json:"verses,omitempty"`
}

func (v Verse) Contains(posMs int) bool {
	return posMs >= v.StartMs && posMs < v.EndMs
}

func (s *Section) VerseAt(posMs int) *Verse {
	for i := range s.Verses {
		if s.Verses[i].Contains(posMs) {
			return &s.Verses[i]
		}
	}
	return nil
}

type Liturgy struct {
	ID             int64  `json:"id"`
	Slug           string `json:"slug"`
	Name           string `json:"name"`
	NameAm         string `json:"name_am"`
	ContentVersion int64  `json:"content_version"`
}

type LiturgyContent struct {
	Liturgy  Liturgy   `json:"liturgy"`
	Sections []Section `json:"sections"`
}
