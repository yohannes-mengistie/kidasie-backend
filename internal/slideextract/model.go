package slideextract

type Draft struct {
	SchemaVersion int         `json:"schema_version"`
	SourcePDF     string      `json:"source_pdf"`
	Audio         *AudioAsset `json:"audio,omitempty"`
	Pages         []Page      `json:"pages"`
}

type AudioAsset struct {
	FileName  string `json:"file_name"`
	SizeBytes int64  `json:"size_bytes"`
	MIMEType  string `json:"mime_type"`
	SHA256    string `json:"sha256"`
}

type Page struct {
	Number            int    `json:"number"`
	Kind              string `json:"kind"`
	RoleSuggestion    string `json:"role_suggestion,omitempty"`
	TextEthiopicOCR   string `json:"text_ethiopic_ocr,omitempty"`
	TextEnglishOCR    string `json:"text_english_ocr,omitempty"`
	NeedsReview       bool   `json:"needs_review"`
	ExtractionWarning string `json:"extraction_warning,omitempty"`
}

const (
	KindContent    = "content"
	KindRoleHeader = "role_header"
	KindMixed      = "mixed"
	KindEmpty      = "empty"
)
