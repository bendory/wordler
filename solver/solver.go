package solver

import (
	"fmt"
	"regexp"

	"wordler/wordlist"
)

const (
	CORRECT   = '+'
	ELSEWHERE = '*'
	NIL       = '_'
)

var verbose = false

// Solver is a wordle guesser.
type Solver struct {
	have map[byte]bool // letters that we know we have
	w    *wordlist.WordList
}

// Option represents a constraint to place on the Dictionary.
type Option interface {
	apply(*Solver)
}

// KeepOnlyOption specifies that Solver should only include words that match
// the given expression.
type KeepOnlyOption struct {
	Exp *regexp.Regexp
}

// apply fulfills the Option interface
func (k KeepOnlyOption) apply(s *Solver) {
	s.w.KeepOnly(k.Exp)
}

// DeleteOption specifies that Solver should exclude words that match
// the given expression.
type DeleteOption struct {
	Exp *regexp.Regexp
}

// apply fulfills the Option interface
func (d DeleteOption) apply(s *Solver) {
	s.w.Delete(d.Exp)
}

// New returns a new Solver populated with a Dictionary.
func New(options ...Option) (*Solver, error) {
	var (
		g   *Solver
		w   *wordlist.WordList
		err error
	)

	if w, err = wordlist.NewDictionary(); err == nil {
		g = &Solver{
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
func (s *Solver) Guess() string {
	return s.w.Random()
}

// React "reacts" to the scored guess by filtering out excluded words from our
// WordList.
func (s *Solver) React(guess, response string) {
	if len(guess) != len(response) {
		panic(fmt.Sprintf("guess len(%v)==%d; response len(%v) == %d", guess, len(guess), response, len(response)))
	}
	if s.have == nil {
		s.have = make(map[byte]bool)
	}

	matches := 0
	keepOnly := "^" // letters in required positions

	for i, r := range response {
		c := guess[i]

		switch r {
		case CORRECT:
			matches++
			keepOnly += string(c)
			s.have[c] = true

		case ELSEWHERE:
			keepOnly += "[^" + string(c) + "]"
			s.w.KeepOnly(regexp.MustCompile(string(c)))
			s.have[c] = true

		case NIL:
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
}

// Remaining returns the number of word remaining for guessing.
func (s *Solver) Remaining() int {
	if s == nil || s.w == nil {
		return 0
	}
	return s.w.Length()
}

func debug(f string, args ...interface{}) {
	if verbose {
		fmt.Printf(f, args...)
		fmt.Println()
	}
}
