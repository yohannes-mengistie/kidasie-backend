// Package liturgydoc decodes the normalized liturgy sources under
// content/updated and joins the shared beginning (ሥርዓተ ቅዳሴ) with a single
// anaphora into one importable document.
package liturgydoc

// Entry is one numbered utterance in a normalized liturgy source. Consecutive
// entries that share a Number belong to the same source group; the source
// numbering is advisory only, because pagination regroups entries to a size
// budget before import.
type Entry struct {
	Number         string `json:"number"`
	Role           string `json:"role"`
	Subtype        string `json:"subtype,omitempty"`
	TextGeez       string `json:"text_geez,omitempty"`
	TextAmharic    string `json:"text_amharic,omitempty"`
	TextEnglish    string `json:"text_english,omitempty"`
	RubricGeez     string `json:"rubric_geez,omitempty"`
	RubricAmharic  string `json:"rubric_amharic,omitempty"`
	RubricEnglish  string `json:"rubric_english,omitempty"`
	OCRNote        string `json:"ocr_note,omitempty"`
	CrossReference string `json:"cross_reference,omitempty"`
	NeedsReview    *bool  `json:"needs_review,omitempty"`

	// number resolved to an integer by Decode. Ge'ez numerals such as ፻፸፭
	// and plain decimals such as 175 both resolve to 175.
	group int
}

// Group reports the source group number the entry was decoded with.
func (e Entry) Group() int {
	return e.group
}

// Document is a decoded liturgy source file. Title is absent in the shared
// beginning and present in every anaphora.
type Document struct {
	Title   string  `json:"title,omitempty"`
	Entries []Entry `json:"entries"`
}

// Options names the liturgy the joined document describes.
type Options struct {
	Slug           string
	Name           string
	NameAm         string
	SectionTitle   string
	SectionTitleAm string

	// Pagination overrides Budget's defaults when non-zero.
	Pagination Budget
}

// Stats reports what a conversion produced, for the operator running it.
type Stats struct {
	BeginningEntries int
	AnaphoraEntries  int
	Pages            int
	SourceGroups     int
	Verses           int
	ReviewSegments   int
	OversizePages    int
	LargestPageRunes int
}
