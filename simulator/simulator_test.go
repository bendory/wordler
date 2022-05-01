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

// Load implements wordlist.DictionaryLoader.Load() by generating a list of word
// permutations.
func (f *fakeLoader) Load(_ ...wordlist.Option) (*wordlist.WordList, error) {
	if f.length < 2 || f.length > 5 {
		return nil, fmt.Errorf("cannot make fake loader for wordlength %d", f.length)
	}

	// Generate a list of f.length letters
	var list []string
	for i := 0; i < f.length; i++ {
		list = append(list, fmt.Sprintf("%c", 'a'+i))
	}
	letters := append([]string{}, list...)

	// Generate every possible "word" of length f.length from the list of
	// letters.
	for len(list[0]) < f.length {
		var next []string
		for i := 0; i < len(list); i++ {
			for _, l := range letters {
				next = append(next, list[i]+l)
			}
		}
		list = next
	}

	f.words = list
	return wordlist.New(f.words), nil
}

func TestSimulations(t *testing.T) {
	args := &puzzler.Args{Hard: true, Guesses: 1, Dictionary: puzzler.LocalDictionary}

	was := wordlist.Loader
	f := &fakeLoader{}
	wordlist.Loader = f
	defer func() { wordlist.Loader = was }()

	// Runtime is O((l^l)^2) -- so running simulations with wordlists >4
	// takes... forever. Length 5 --> 9MM+ test cases, and I don't think we get
	// any extra benefit from it.
	for l := 2; l <= 4; l++ {
		f.length = l
		t.Run(fmt.Sprint(l), func(t *testing.T) {
			if _, err := f.Load(); err != nil {
				t.Fatal(err)
			}
			for _, solution := range f.words {
				t.Run(solution, func(t *testing.T) {
					args.Solution = solution
					for _, guess := range f.words {
						t.Run(guess, func(t *testing.T) {
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
							if p.Words() == 0 {
								t.Error("no words left.")
							}
							wordsBefore := p.Words()

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
							if p.Words() == 0 {
								t.Error("no words left.")
							}
							if wordsBefore <= p.Words() {
								t.Errorf("%d words before <= %d words after", wordsBefore, p.Words())
							}
						})
					}
				})
			}
		})
	}
}
