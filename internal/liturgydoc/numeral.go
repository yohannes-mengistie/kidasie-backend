package liturgydoc

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	geezHundred     = '፻'
	geezTenThousand = '፼'
)

// geezDigits holds the Ethiopic unit and tens digits. The multipliers ፻ and ፼
// are handled separately because they scale the group before them.
var geezDigits = map[rune]int{
	'፩': 1, '፪': 2, '፫': 3, '፬': 4, '፭': 5,
	'፮': 6, '፯': 7, '፰': 8, '፱': 9,
	'፲': 10, '፳': 20, '፴': 30, '፵': 40, '፶': 50,
	'፷': 60, '፸': 70, '፹': 80, '፺': 90,
}

// parseNumber resolves a source entry number, which is written either as a
// plain decimal (the shared beginning) or as a Ge'ez numeral (the anaphoras).
func parseNumber(value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, fmt.Errorf("number is required")
	}

	if number, err := strconv.Atoi(value); err == nil {
		if number <= 0 {
			return 0, fmt.Errorf("number %q must be positive", value)
		}

		return number, nil
	}

	number, err := parseGeezNumeral(value)
	if err != nil {
		return 0, fmt.Errorf("number %q: %w", value, err)
	}

	return number, nil
}

// parseGeezNumeral evaluates an Ethiopic numeral. Ge'ez is positional by
// multiplier rather than by place: the group before ፼ is scaled by ten
// thousand and the group before ፻ by one hundred, so ፻፸፭ is 100+70+5 and
// ፫፻፲፰ is 3*100+10+8. A bare multiplier carries an implicit leading ፩.
func parseGeezNumeral(value string) (int, error) {
	runes := []rune(value)

	if index := lastIndexRune(runes, geezTenThousand); index >= 0 {
		high, err := parseGeezSegment(runes[:index], 1)
		if err != nil {
			return 0, err
		}

		low, err := parseGeezNumeral(string(runes[index+1:]))
		if err != nil && len(runes[index+1:]) > 0 {
			return 0, err
		}

		return high*10000 + low, nil
	}

	if index := lastIndexRune(runes, geezHundred); index >= 0 {
		high, err := parseGeezSegment(runes[:index], 1)
		if err != nil {
			return 0, err
		}

		low, err := parseGeezSegment(runes[index+1:], 0)
		if err != nil {
			return 0, err
		}

		return high*100 + low, nil
	}

	return parseGeezSegment(runes, 0)
}

// parseGeezSegment sums the plain digits in a multiplier-free segment,
// returning fallback when the segment is empty.
func parseGeezSegment(runes []rune, fallback int) (int, error) {
	if len(runes) == 0 {
		return fallback, nil
	}

	total := 0
	previous := 0

	for _, symbol := range runes {
		digit, known := geezDigits[symbol]
		if !known {
			return 0, fmt.Errorf("unsupported numeral %q", string(symbol))
		}

		if previous != 0 && digit >= previous {
			return 0, fmt.Errorf(
				"numeral digits must descend, found %q after %q",
				string(symbol),
				geezSymbol(previous),
			)
		}

		total += digit
		previous = digit
	}

	return total, nil
}

func lastIndexRune(runes []rune, target rune) int {
	for index := len(runes) - 1; index >= 0; index-- {
		if runes[index] == target {
			return index
		}
	}

	return -1
}

func geezSymbol(digit int) string {
	for symbol, value := range geezDigits {
		if value == digit {
			return string(symbol)
		}
	}

	return strconv.Itoa(digit)
}
