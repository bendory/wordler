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
	have map[byte]bool      // letters that we know we have
	s    *wordlist.WordList // words that are valid solutions
	g    *wordlist.WordList // words that are valid guesses
}

// From returns a new Solver created from the given list of words.
func From(dictionary []string) *Solver {
	return &Solver{
		have: make(map[byte]bool, 26),
		s:    wordlist.New(dictionary),
		g:    wordlist.New(dictionary),
	}
}

// New returns a new Solver populated with the default Dictionary.
func New(options ...wordlist.Option) (*Solver, error) {
	var (
		s   *Solver
		w   *wordlist.WordList
		err error
	)

	if w, err = wordlist.NewDictionary(options...); err == nil {
		s = &Solver{
			have: make(map[byte]bool, 26),
			s:    w,
			g:    w.Clone(),
		}
	}

	return s, err
}

// Guess provides a guess from remaining words
func (s *Solver) Guess() string {
	return s.s.OptimalGuess()
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
				// s.s will be filtered by keepOnly at end of loop.
				// s.g is filtered here because keepOnly is more strict than
				// hard mode.
				keepOnly[i] = string(c)
				s.have[c] = true
				s.g.KeepOnly(regexp.MustCompile("^" + strings.Repeat(".", i) + string(c)))

			case wordler.ELSEWHERE:
				// s.s will be filtered by keepOnly at end of loop; s.g won't.
				keepOnly[i] = "[^" + string(c) + "]"
				r := regexp.MustCompile(string(c))
				s.s.KeepOnly(r) // solution must contain c
				s.g.KeepOnly(r) // future guesses must contain c
				s.have[c] = true

			case wordler.NIL:
				keepOnly[i] = "[^" + string(c) + "]"
				if !s.have[c] {
					s.s.Delete(regexp.MustCompile(string(c)))
					// No update to s.g -- can still use this letter in guesses.
				}
			}
		}
	}

	debug("found %d matches", matches)
	if matches == len(guess) {
		// complete match!
		s.s = wordlist.New([]string{guess})
		s.g = s.s.Clone()
	} else {
		p := "^" + strings.Join(keepOnly, "") + "$"
		debug("keepOnly: '%v'", p)
		s.s.KeepOnly(regexp.MustCompile(p))
		// always eliminate guess itself
		s.g.Delete(regexp.MustCompile(guess))
	}
	return nil
}

// Remaining returns the number of possible solutions remaining.
func (s *Solver) Remaining() int {
	if s == nil || s.s == nil {
		return 0
	}
	return s.s.Length()
}

// NotInWordle is used to report that the word is not found in the wordle
// dictionary; the word is removed from our list of remaining entries.
func (s *Solver) NotInWordle(not string) {
	if s == nil || s.s == nil {
		return
	}
	r := regexp.MustCompile("^" + not + "$")
	s.s.Delete(r)
	s.g.Delete(r)
}

// debug prints debug logs
func debug(f string, args ...interface{}) {
	if verbose {
		fmt.Printf(f, args...)
		fmt.Println()
	}
}
