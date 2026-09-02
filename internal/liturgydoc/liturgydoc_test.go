package liturgydoc

import (
	"strings"
	"testing"
)

func TestParseNumberReadsGeezAndDecimal(t *testing.T) {
	cases := map[string]int{
		"1":    1,
		"193":  193,
		"፩":    1,
		"፲":    10,
		"፲፰":   18,
		"፸፮":   76,
		"፻":    100,
		"፻፸፭":  175,
		"፫፻፲፰": 318,
		"፼":    10000,
	}

	for value, want := range cases {
		got, err := parseNumber(value)
		if err != nil {
			t.Fatalf("parseNumber(%q) returned error: %v", value, err)
		}

		if got != want {
			t.Errorf("parseNumber(%q) = %d, want %d", value, got, want)
		}
	}
}

func TestParseNumberRejectsBadInput(t *testing.T) {
	for _, value := range []string{"", "0", "-3", "abc", "፩፲", "x፻"} {
		if _, err := parseNumber(value); err == nil {
			t.Errorf("parseNumber(%q) accepted an invalid number", value)
		}
	}
}

func TestDecodeAcceptsObjectAndArray(t *testing.T) {
	object := `{"title":"ቅዳሴ","entries":[
		{"number":"፩","role":"priest","text_geez":"ሀ"},
		{"number":"፪","role":"people","text_amharic":"ለ"}]}`

	document, err := Decode(strings.NewReader(object))
	if err != nil {
		t.Fatalf("Decode object returned error: %v", err)
	}

	if document.Title != "ቅዳሴ" {
		t.Errorf("title = %q, want %q", document.Title, "ቅዳሴ")
	}

	if got := document.Entries[1].Group(); got != 2 {
		t.Errorf("second entry group = %d, want 2", got)
	}

	array := `[{"number":"1","role":"deacon","text_geez":"ሀ"}]`

	if _, err := Decode(strings.NewReader(array)); err != nil {
		t.Fatalf("Decode array returned error: %v", err)
	}
}

func TestDecodeRejectsNonSequentialNumbers(t *testing.T) {
	source := `{"entries":[
		{"number":"፩","role":"priest","text_geez":"ሀ"},
		{"number":"፫","role":"priest","text_geez":"ለ"}]}`

	_, err := Decode(strings.NewReader(source))
	if err == nil {
		t.Fatal("Decode accepted a gap in the entry numbering")
	}

	if !strings.Contains(err.Error(), "expected 1 or 2") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDecodeRejectsEmptyEntry(t *testing.T) {
	source := `{"entries":[{"number":"1","role":"priest"}]}`

	if _, err := Decode(strings.NewReader(source)); err == nil {
		t.Fatal("Decode accepted an entry with no text")
	}
}

func TestPaginatePacksSmallEntriesAndKeepsEntriesWhole(t *testing.T) {
	entries := []Entry{
		{Number: "1", Role: "priest", TextGeez: strings.Repeat("ሀ", 60)},
		{Number: "2", Role: "people", TextGeez: strings.Repeat("ለ", 60)},
		{Number: "3", Role: "priest", TextGeez: strings.Repeat("ሐ", 60)},
		{Number: "4", Role: "priest", TextGeez: strings.Repeat("መ", 250)},
	}

	pages := paginate(entries, Budget{TargetRunes: 200}, 1)

	if len(pages) != 2 {
		t.Fatalf("page count = %d, want 2", len(pages))
	}

	if len(pages[0].Entries) != 3 {
		t.Errorf(
			"first page holds %d entries, want the 3 that fit the budget",
			len(pages[0].Entries),
		)
	}

	// The final entry is longer than the whole budget. It must still be a
	// single page, because splitting it would divide one prayer's Ge'ez.
	if len(pages[1].Entries) != 1 || pages[1].Runes() != 250 {
		t.Errorf(
			"oversized entry was split: %d entries, %d runes",
			len(pages[1].Entries),
			pages[1].Runes(),
		)
	}

	for index, page := range pages {
		if page.Number != index+1 {
			t.Errorf("page %d numbered %d", index, page.Number)
		}
	}
}

func TestPaginateStartsPageAtPrayerHeader(t *testing.T) {
	entries := []Entry{
		{Number: "1", Role: "priest", TextGeez: "ሀ"},
		{
			Number:   "2",
			Role:     "rubric",
			Subtype:  "prayer_header",
			TextGeez: "ጸሎተ ፍትሐት ።",
		},
		{Number: "2", Role: "priest", TextGeez: "ለ"},
	}

	pages := paginate(entries, Budget{TargetRunes: 1000}, 1)

	if len(pages) != 2 {
		t.Fatalf("page count = %d, want 2", len(pages))
	}

	if pages[1].Entries[0].Subtype != "prayer_header" {
		t.Error("prayer header did not open the second page")
	}
}

func newTestDocument(numbers []string, text string) *Document {
	entries := make([]Entry, 0, len(numbers))
	for _, number := range numbers {
		entries = append(entries, Entry{
			Number:   number,
			Role:     "priest",
			TextGeez: text,
		})
	}

	document := &Document{Entries: entries}
	if err := resolveGroups(document.Entries); err != nil {
		panic(err)
	}

	return document
}

func testOptions() Options {
	return Options{
		Slug:           "test-liturgy",
		Name:           "Test Liturgy",
		NameAm:         "የሙከራ ቅዳሴ",
		SectionTitle:   "Complete Test Liturgy",
		SectionTitleAm: "ሙሉ የሙከራ ቅዳሴ",
		Pagination:     Budget{TargetRunes: 100},
	}
}

func TestConvertNumbersPagesContinuouslyAcrossTheJoin(t *testing.T) {
	beginning := newTestDocument(
		[]string{"1", "2", "3", "4"},
		strings.Repeat("ሀ", 90),
	)
	anaphora := newTestDocument(
		[]string{"፩", "፪"},
		strings.Repeat("ለ", 90),
	)

	document, stats, err := Convert(beginning, anaphora, testOptions())
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if stats.Pages != 6 {
		t.Fatalf("pages = %d, want 6", stats.Pages)
	}

	verses := document.Sections[0].Verses
	if len(verses) != 6 {
		t.Fatalf("verses = %d, want 6", len(verses))
	}

	// Both sources restart their own numbering at one; the joined pages must
	// not.
	for index, verse := range verses {
		if verse.SourcePage == nil {
			t.Fatalf("verses[%d] has no page", index)
		}

		if *verse.SourcePage != index+1 {
			t.Errorf(
				"verses[%d] on page %d, want %d",
				index,
				*verse.SourcePage,
				index+1,
			)
		}
	}

	if got := verses[3].SourceKind; got != kindBeginning {
		t.Errorf("last beginning verse kind = %q", got)
	}

	if got := verses[4].SourceKind; got != kindAnaphora {
		t.Errorf("first anaphora verse kind = %q", got)
	}
}

func TestConvertStartsAnaphoraOnItsOwnPage(t *testing.T) {
	// Both documents are far below the budget, so a single paginator would
	// merge the join onto one page.
	beginning := newTestDocument([]string{"1"}, "ሀ")
	anaphora := newTestDocument([]string{"፩"}, "ለ")

	document, stats, err := Convert(beginning, anaphora, testOptions())
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if stats.Pages != 2 {
		t.Errorf("pages = %d, want the anaphora to open a new page", stats.Pages)
	}

	verses := document.Sections[0].Verses
	if *verses[0].SourcePage == *verses[1].SourcePage {
		t.Error("anaphora shares a page with the beginning")
	}
}

func TestConvertEmitsRubricBeforeItsUtteranceOnTheSamePage(t *testing.T) {
	beginning := &Document{Entries: []Entry{{
		Number:        "1",
		Role:          "assistant_priest",
		RubricGeez:    "ቡራኬ ላዕለ ሕዝብ ።",
		RubricAmharic: "ሕዝቡን ይባርክ ።",
		RubricEnglish: "Let him bless the people.",
		TextGeez:      "ኦ ሥሉስ ቅዱስ ።",
		TextAmharic:   "ልዩ ሦስት ሆይ ።",
		TextEnglish:   "O Holy Trinity.",
		OCRNote:       "printed twice",
	}}}
	if err := resolveGroups(beginning.Entries); err != nil {
		t.Fatalf("resolveGroups: %v", err)
	}

	anaphora := newTestDocument([]string{"፩"}, "ለ")

	document, _, err := Convert(beginning, anaphora, testOptions())
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	verses := document.Sections[0].Verses
	if verses[0].Role != "rubric" {
		t.Fatalf("first verse role = %q, want rubric", verses[0].Role)
	}

	if verses[0].TextEn != "Let him bless the people." {
		t.Errorf("rubric English = %q", verses[0].TextEn)
	}

	if verses[1].Role != "assistant_priest" {
		t.Errorf("second verse role = %q", verses[1].Role)
	}

	if *verses[0].SourcePage != *verses[1].SourcePage {
		t.Error("rubric was separated from its utterance")
	}

	if verses[1].SourceNote != "printed twice" {
		t.Errorf("source note = %q", verses[1].SourceNote)
	}
}

func TestConvertTreatsAbsentNeedsReviewAsReviewed(t *testing.T) {
	flagged := true
	beginning := &Document{Entries: []Entry{
		{Number: "1", Role: "priest", TextGeez: "ሀ"},
		{Number: "2", Role: "people", TextGeez: "ለ", NeedsReview: &flagged},
	}}
	if err := resolveGroups(beginning.Entries); err != nil {
		t.Fatalf("resolveGroups: %v", err)
	}

	anaphora := newTestDocument([]string{"፩"}, "ሐ")

	document, stats, err := Convert(beginning, anaphora, testOptions())
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if stats.ReviewSegments != 1 {
		t.Errorf("review segments = %d, want 1", stats.ReviewSegments)
	}

	verses := document.Sections[0].Verses
	if verses[0].SourceNeedsReview {
		t.Error("an unflagged entry was marked as needing review")
	}

	if !verses[1].SourceNeedsReview {
		t.Error("a flagged entry lost its review flag")
	}

	// "people" is the source's name for the congregation role.
	if verses[1].Role != "congregation" {
		t.Errorf("role = %q, want congregation", verses[1].Role)
	}
}

func TestConvertRejectsAnUnknownRole(t *testing.T) {
	beginning := &Document{Entries: []Entry{
		{Number: "1", Role: "cantor-in-chief", TextGeez: "ሀ"},
	}}
	if err := resolveGroups(beginning.Entries); err != nil {
		t.Fatalf("resolveGroups: %v", err)
	}

	_, _, err := Convert(beginning, newTestDocument([]string{"፩"}, "ለ"), testOptions())
	if err == nil {
		t.Fatal("Convert accepted an unknown role")
	}
}
