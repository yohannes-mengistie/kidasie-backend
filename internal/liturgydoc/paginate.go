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

	// MinRunes is the fill a page should reach before it is closed. Closing
	// greedily at the target alone strands a short entry on its own page
	// whenever the entry after it is large, which is what makes one page look
	// nothing like the next. Below this floor the next entry is pulled in even
	// though it overshoots the target, as long as it stays under MaxRunes.
	// Defaults to just under half the target.
	MinRunes int

	// MaxRunes is the hard ceiling for that overshoot. An entry that cannot
	// fit under it opens its own page instead. Defaults to half again the
	// target. A single entry longer than MaxRunes still gets a page of its
	// own, since an entry is never split.
	MaxRunes int
}

func (b Budget) target() int {
	if b.TargetRunes > 0 {
		return b.TargetRunes
	}

	return DefaultTargetRunes
}

func (b Budget) minimum() int {
	if b.MinRunes > 0 {
		return b.MinRunes
	}

	return b.target() * 45 / 100
}

func (b Budget) maximum() int {
	if b.MaxRunes > 0 {
		return b.MaxRunes
	}

	return b.target() * 3 / 2
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

		if len(current) > 0 && breaksBefore(current, currentRunes, entry, size, budget) {
			flush()
		}

		current = append(current, entry)
		currentRunes += size
	}

	flush()

	return mergeShortPages(pages, budget, startNumber)
}

// mergeShortPages folds a page that still fell below the floor back into the
// page before it. The forward pass cannot do this on its own: when a short
// entry is followed by one too large to join it, the short entry has already
// opened a page by the time that is known, and only a look backwards can put
// it where it belongs.
func mergeShortPages(pages []Page, budget Budget, startNumber int) []Page {
	merged := make([]Page, 0, len(pages))

	for _, page := range pages {
		if len(merged) > 0 && joinsPrevious(merged[len(merged)-1], page, budget) {
			previous := &merged[len(merged)-1]
			previous.Entries = append(previous.Entries, page.Entries...)

			continue
		}

		merged = append(merged, page)
	}

	for i := range merged {
		merged[i].Number = startNumber + i
	}

	return merged
}

// joinsPrevious reports whether a page is short enough to fold backwards, and
// whether the page before it has the room.
func joinsPrevious(previous, page Page, budget Budget) bool {
	if startsPage(page.Entries[0]) {
		return false
	}

	if page.Runes() >= budget.minimum() {
		return false
	}

	return previous.Runes()+page.Runes() <= budget.maximum()
}

// breaksBefore reports whether entry opens a new page rather than joining the
// page built so far.
func breaksBefore(current []Entry, currentRunes int, entry Entry, size int, budget Budget) bool {
	// A header is the title of what follows, so it must never be the last
	// thing on a page. Whatever comes after it joins it however long it is.
	if onlyHeaders(current) {
		return false
	}

	if startsPage(entry) {
		return true
	}

	if currentRunes+size <= budget.target() {
		return false
	}

	// Over the target. Take the entry anyway rather than close a page that is
	// barely filled, unless doing so would overshoot far enough to be its own
	// kind of inconsistency.
	return currentRunes >= budget.minimum() || currentRunes+size > budget.maximum()
}

// startsPage reports whether an entry must open a page. A prayer header is
// the title of what follows, so it belongs at the top of a page rather than
// stranded at the bottom of the previous one.
func startsPage(entry Entry) bool {
	return entry.Subtype == "prayer_header"
}

// onlyHeaders reports whether the page built so far is nothing but titles.
func onlyHeaders(entries []Entry) bool {
	for _, entry := range entries {
		if !startsPage(entry) {
			return false
		}
	}

	return true
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
