package simulator

import (
	"fmt"
	"testing"

	"wordler/puzzler"
	"wordler/solver"
	"wordler/wordlist"
)

// fakeLoader implements wordlist.DictionaryLoader.
type fakeLoader struct {
	length int // wordlength
	words  []string
}

func (f *fakeLoader) Load(_ ...wordlist.Option) (*wordlist.WordList, error) {
	switch f.length {
	case 1:
		f.words = []string{"a", "b", "c", "d", "e"}
	case 2:
		f.words = []string{"aa", "ab", "ac", "ba", "bb", "bc", "ca", "cb", "cc"}
	case 3:
		f.words = []string{
			"aaa", "aab", "aac",
			"aba", "abb", "abc",
			"aca", "acb", "acc",
			"baa", "bab", "bac",
			"bba", "bbb", "bbc",
			"bca", "bcb", "bcc",
			"caa", "cab", "cac",
			"cba", "cbb", "cbc",
			"cca", "ccb", "ccc",
		}
	default:
		return nil, fmt.Errorf("cannot make fake loader for wordlength %d", f.length)
	}
	return wordlist.New(f.words), nil
}

func TestSimulations(t *testing.T) {
	args := &puzzler.Args{Hard: true, Guesses: 1}
	for l := 1; l <= 3; l++ {
		t.Run(fmt.Sprint(l), func(t *testing.T) {
			f := &fakeLoader{length: l}
			wordlist.Loader = f
			f.Load()

			for _, solution := range f.words {
				t.Run(solution, func(t *testing.T) {
					for _, guess := range f.words {
						t.Run(guess, func(t *testing.T) {
							args.Solution = solution
							p, err := puzzler.New(args)
							if err != nil {
								t.Fatalf("Failed to make a Puzzler: %v", err)
							}
							s, err := solver.New()
							if err != nil {
								t.Fatalf("Failed to make a Solver: %v", err)
							}

							if p.Words() != s.Remaining() {
								t.Errorf("%d Puzzler words != %d Solver words", p.Words(), s.Remaining())
							}

							response, err := p.Guess(guess)
							if err != nil {
								t.Errorf("ERROR: %v", err)
							}
							if err := s.React(guess, response); err != nil {
								t.Errorf("ERROR: %v", err)
							}

							if p.Words() != s.Remaining() {
								t.Errorf("%d Puzzler words != %d Solver words\n", p.Words(), s.Remaining())
							}
						})
					}
				})
			}
		})
	}
}
