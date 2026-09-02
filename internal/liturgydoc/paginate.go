package liturgydoc

import "strings"

// DefaultTargetRunes is the per-page text budget used when Options does not
// override it. It counts every rune the reader can see on a page, across
// Ge'ez, Amharic and English plus any rubric.
//
// The source groups average ~314 runes and nearly half fall below 300, so
// paginating on them one-to-one produces 270-370 very short pages per
// service. Packing to 1000 runes roughly halves the page count while keeping
// the ninetieth-percentile page near one budget's worth of scrolling.
const DefaultTargetRunes = 1000

// Budget controls how entries are packed onto pages.
type Budget struct {
	// TargetRunes is the soft ceiling for a page. A page is closed before it
	// would exceed the target, so only a single entry longer than the target
	// can produce a page above it.
	TargetRunes int
}

func (b Budget) target() int {
	if b.TargetRunes > 0 {
		return b.TargetRunes
	}

	return DefaultTargetRunes
}

// Page is a run of consecutive entries presented together, numbered
// continuously across the joined document.
type Page struct {
	Number  int
	Entries []Entry
}

// Runes reports the visible length of the page.
func (p Page) Runes() int {
	total := 0
	for _, entry := range p.Entries {
		total += entryRunes(entry)
	}

	return total
}

// paginate packs entries onto pages without ever splitting an entry, so the
// Ge'ez and Amharic of a single utterance are never divided across a page
// break. Pages start at startNumber and increase by one.
func paginate(entries []Entry, budget Budget, startNumber int) []Page {
	target := budget.target()

	pages := make([]Page, 0, len(entries))
	current := make([]Entry, 0, 8)
	currentRunes := 0

	flush := func() {
		if len(current) == 0 {
			return
		}

		pages = append(pages, Page{
			Number:  startNumber + len(pages),
			Entries: current,
		})
		current = make([]Entry, 0, 8)
		currentRunes = 0
	}

	for _, entry := range entries {
		size := entryRunes(entry)

		if len(current) > 0 &&
			(startsPage(entry) || currentRunes+size > target) {
			flush()
		}

		current = append(current, entry)
		currentRunes += size
	}

	flush()

	return pages
}

// startsPage reports whether an entry must open a page. A prayer header is
// the title of what follows, so it belongs at the top of a page rather than
// stranded at the bottom of the previous one.
func startsPage(entry Entry) bool {
	return entry.Subtype == "prayer_header"
}

func entryRunes(entry Entry) int {
	return countRunes(entry.TextGeez) +
		countRunes(entry.TextAmharic) +
		countRunes(entry.TextEnglish) +
		countRunes(entry.RubricGeez) +
		countRunes(entry.RubricAmharic) +
		countRunes(entry.RubricEnglish)
}

func countRunes(value string) int {
	return len([]rune(strings.TrimSpace(value)))
}
