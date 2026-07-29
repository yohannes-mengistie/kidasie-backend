package domain

type Verse struct {
	ID int64 `json:"id"`
	Order int `json:"order"`
	TextGeez string `json:"text_geez"`
	TextAm string `json:"text_am"`
	TextEn string `json:"text_en"`
	Chanter string `json:"chanter"`
	StartMs int `json:"start_ms"`
	EndMs int	`json:"end_ms"`

}

func (v Verse) Contains(posMs int) bool {
	return posMs >= v.StartMs && posMs <= v.EndMs
}

type Section struct {
	ID int64 `json:"id"`
	Title string `json:"title"`
	TitleAm string `json:"title_am"`
	AudioURL string `json:"audio_url"`
	Verses []Verse `json:"verses"`
}

func (s *Section) VerseAt(posMs int) *Verse {
	for _, v := range s.Verses {
		if v.Contains(posMs) {
			return &v
		}
	}
	return nil
}

type Liturgy struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	NameAm string `json:"name_am"`
}