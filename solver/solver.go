package solver

import (
	"fmt"
	"regexp"
	"strings"

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
	pattern := fmt.Sprintf("^[%s]{%d}$", []byte{wordler.CORRECT, wordler.ELSEWHERE, wordler.NIL}, len(guess))
	if r := regexp.MustCompile(pattern); !r.MatchString(response) {
		return fmt.Errorf("invalid response: response must match %#v", pattern)
	}
	if s.have == nil {
		s.have = make(map[byte]bool)
	}

	matches := 0
	keepOnly := make([]string, len(guess)) // letters in required positions

	// Need to process response signals in this order; see test case for
	// combination of guess "carer" for word "foyer" in solver_test to
	// understand why.
	for _, reaction := range []rune{wordler.CORRECT, wordler.ELSEWHERE, wordler.NIL} {
		for i, r := range response {
			if r != reaction {
				continue
			}
			c := guess[i]

			switch r {
			case wordler.CORRECT:
				matches++
				keepOnly[i] = string(c)
				s.have[c] = true

			case wordler.ELSEWHERE:
				keepOnly[i] = "[^" + string(c) + "]"
				s.w.KeepOnly(regexp.MustCompile(string(c)))
				s.have[c] = true

			case wordler.NIL:
				keepOnly[i] = "[^" + string(c) + "]"
				if !s.have[c] {
					s.w.Delete(regexp.MustCompile(string(c)))
				}
			}
		}
	}

	debug("found %d matches", matches)
	if matches == len(guess) {
		// complete match!
		s.w = wordlist.New([]string{guess})
	} else {
		p := "^" + strings.Join(keepOnly, "") + "$"
		debug("keepOnly: '%v'", p)
		s.w.KeepOnly(regexp.MustCompile(p))
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
