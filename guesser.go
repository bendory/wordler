package wordler

import (
	"regexp"
	"strings"
)

const (
	RIGHT_LETTER_RIGHT_PLACE = '+'
	RIGHT_LETTER_WRONG_PLACE = '*'
	LETTER_NOT_IN_WORD       = ' '
)

type Guesser struct {
	w *WordList
}

// NewGuesser returns a new Guesser populated with a Dictionary.
func NewGuesser() (*Guesser, error) {
	var (
		g   *Guesser
		w   *WordList
		err error
	)

	if w, err = NewDictionary(); err == nil {
		g = &Guesser{w: w}
	}
	return g, err
}

// Guess provides a random guess from remaining words
func (g *Guesser) Guess() string {
	// map iteration goes in a random order!
	if g.w != nil {
		for guess, _ := range g.w.words {
			return guess
		}
	}
	return ""
}

// React "reacts" to the scored guess by filtering out excluded words from our
// WordList.
func (g *Guesser) React(guess, response string) {
	matches := 0
	matchFilter := "^"       // letters in required spaces
	missFilter := ""         // excluded letters
	hasFilters := []string{} // letters included in some other space

	for i, r := range response {
		c := guess[i]

		switch r {
		case RIGHT_LETTER_RIGHT_PLACE:
			matches++
			matchFilter += string(c)

		case RIGHT_LETTER_WRONG_PLACE:
			matchFilter += "."
			hasFilters = append(hasFilters, "^"+strings.Repeat(".", i)+string(c))

		case LETTER_NOT_IN_WORD:
			matchFilter += "."
			missFilter += string(c)
		}

		switch matches {
		case 0:
			// do nothing; no matches
		case len(guess):
			// complete match!
			g.w = NewWordList([]string{guess})
		default:
			// we found some matches, but not a complete match
			matchFilter += "$"
			g.w.KeepOnly(regexp.MustCompile(matchFilter))
		}
	}

	if len(missFilter) > 0 {
		g.w.Delete(regexp.MustCompile("[" + missFilter + "]"))
	}

	if len(hasFilters) > 0 {
		g.w.Delete(regexp.MustCompile(strings.Join(hasFilters, "|")))
	}

	return
}
