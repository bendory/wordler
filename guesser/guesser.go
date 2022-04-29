package guesser

import (
	"regexp"
	"strings"

	"wordler/wordlist"
)

const (
	RIGHT_LETTER_RIGHT_PLACE = '+'
	RIGHT_LETTER_WRONG_PLACE = '*'
	LETTER_NOT_IN_WORD       = ' '
)

type Guesser struct {
	w *wordlist.WordList
}

// New returns a new Guesser populated with a Dictionary.
func New() (*Guesser, error) {
	var (
		g   *Guesser
		w   *wordlist.WordList
		err error
	)

	if w, err = wordlist.NewDictionary(); err == nil {
		g = &Guesser{w: w}
	}
	return g, err
}

// Guess provides a random guess from remaining words
func (g *Guesser) Guess() string {
	return g.w.Random()
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
	}

	switch matches {
	case 0:
		// do nothing; no matches
	case len(guess):
		// complete match!
		g.w = wordlist.New([]string{guess})
	default:
		// we found some matches, but not a complete match
		matchFilter += "$"
		g.w.KeepOnly(regexp.MustCompile(matchFilter))
	}

	if len(missFilter) > 0 {
		g.w.Delete(regexp.MustCompile("[" + missFilter + "]"))
	}

	if len(hasFilters) > 0 {
		g.w.Delete(regexp.MustCompile(strings.Join(hasFilters, "|")))
	}

	return
}

// Remaining returns the number of word remaining for guessing.
func (g *Guesser) Remaining() int {
	if g == nil || g.w == nil {
		return 0
	}
	return g.w.Length()
}
