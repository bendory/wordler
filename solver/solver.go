package solver

import (
	"fmt"
	"regexp"

	"wordler"
	"wordler/wordlist"
)

var verbose = false

// Solver is a wordle guesser.
type Solver struct {
	have map[byte]bool // letters that we know we have
	w    *wordlist.WordList
}

// New returns a new Solver populated with a Dictionary.
func New(options ...wordlist.Option) (*Solver, error) {
	var (
		g   *Solver
		w   *wordlist.WordList
		err error
	)

	if w, err = wordlist.NewDictionary(options...); err == nil {
		g = &Solver{
			have: make(map[byte]bool),
			w:    w,
		}
	}

	return g, err
}

// Guess provides a random guess from remaining words
func (s *Solver) Guess() string {
	// TODO: make this a smart algorithm:
	// - guess most common letters first
	// - avoid guessing double letters
	return s.w.Random()
}

// React "reacts" to the scored guess by filtering out excluded words from our
// WordList.
func (s *Solver) React(guess, response string) error {
	if len(guess) != len(response) {
		return fmt.Errorf("response len(%v)==%d does not match guess len(%v)==%d", response, len(response), guess, len(guess))
	}
	if s.have == nil {
		s.have = make(map[byte]bool)
	}

	matches := 0
	keepOnly := "^" // letters in required positions

	for i, r := range response {
		c := guess[i]

		switch r {
		case wordler.CORRECT:
			matches++
			keepOnly += string(c)
			s.have[c] = true

		case wordler.ELSEWHERE:
			keepOnly += "[^" + string(c) + "]"
			s.w.KeepOnly(regexp.MustCompile(string(c)))
			s.have[c] = true

		case wordler.NIL:
			keepOnly += "[^" + string(c) + "]"
			if !s.have[c] {
				s.w.Delete(regexp.MustCompile(string(c)))
			}
		}
	}

	debug("found %d matches", matches)
	if matches == len(guess) {
		// complete match!
		s.w = wordlist.New([]string{guess})
	} else {
		keepOnly += "$"
		debug("keepOnly: '%v'", keepOnly)
		s.w.KeepOnly(regexp.MustCompile(keepOnly))
	}
	return nil
}

// Remaining returns the number of word remaining for guessing.
func (s *Solver) Remaining() int {
	if s == nil || s.w == nil {
		return 0
	}
	return s.w.Length()
}

// NotInWordle is used to report that the word is not found in the wordle
// dictionary; the word is removed from our list of remaining entries.
func (s *Solver) NotInWordle(not string) {
	if s == nil || s.w == nil {
		return
	}
	s.w.Delete(regexp.MustCompile("^" + not + "$"))
}

// debug prints debug logs
func debug(f string, args ...interface{}) {
	if verbose {
		fmt.Printf(f, args...)
		fmt.Println()
	}
}
