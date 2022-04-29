package guesser

import (
	"fmt"
	"regexp"

	"wordler/wordlist"
)

const (
	RIGHT_LETTER_RIGHT_PLACE = '+'
	RIGHT_LETTER_WRONG_PLACE = '*'
	LETTER_NOT_IN_WORD       = '_'
)

var verbose = false

// Guesser is a wordle guesser.
type Guesser struct {
	have map[byte]bool // letters that we know we have
	w    *wordlist.WordList
}

// Option represents a constraint to place on the Dictionary.
type Option interface {
	apply(*Guesser)
}

// KeepOnlyOption specifies that Guesser should only include words that match
// the given expression.
type KeepOnlyOption struct {
	Exp *regexp.Regexp
}

// apply fulfills the Option interface
func (k KeepOnlyOption) apply(g *Guesser) {
	g.w.KeepOnly(k.Exp)
}

// DeleteOption specifies that Guesser should exclude words that match
// the given expression.
type DeleteOption struct {
	Exp *regexp.Regexp
}

// apply fulfills the Option interface
func (d DeleteOption) apply(g *Guesser) {
	g.w.Delete(d.Exp)
}

// New returns a new Guesser populated with a Dictionary.
func New(options ...Option) (*Guesser, error) {
	var (
		g   *Guesser
		w   *wordlist.WordList
		err error
	)

	if w, err = wordlist.NewDictionary(); err == nil {
		g = &Guesser{
			have: make(map[byte]bool),
			w:    w,
		}
	}

	for _, o := range options {
		o.apply(g)
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
	if len(guess) != len(response) {
		panic(fmt.Sprintf("guess len(%v)==%d; response len(%v) == %d", guess, len(guess), response, len(response)))
	}
	if g.have == nil {
		g.have = make(map[byte]bool)
	}

	matches := 0
	keepOnly := "^" // letters in required positions

	for i, r := range response {
		c := guess[i]

		switch r {
		case RIGHT_LETTER_RIGHT_PLACE:
			matches++
			keepOnly += string(c)
			g.have[c] = true

		case RIGHT_LETTER_WRONG_PLACE:
			keepOnly += "[^" + string(c) + "]"
			g.w.KeepOnly(regexp.MustCompile(string(c)))
			g.have[c] = true

		case LETTER_NOT_IN_WORD:
			keepOnly += "."
			if !g.have[c] {
				g.w.Delete(regexp.MustCompile(string(c)))
			}
		}
	}

	debug("found %d matches", matches)
	if matches == len(guess) {
		// complete match!
		g.w = wordlist.New([]string{guess})
	} else {
		keepOnly += "$"
		debug("keepOnly: '%v'", keepOnly)
		g.w.KeepOnly(regexp.MustCompile(keepOnly))
	}
}

// Remaining returns the number of word remaining for guessing.
func (g *Guesser) Remaining() int {
	if g == nil || g.w == nil {
		return 0
	}
	return g.w.Length()
}

func debug(f string, args ...interface{}) {
	if verbose {
		fmt.Printf(f, args...)
		fmt.Println()
	}
}
