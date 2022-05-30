package solver

import (
	"regexp"
	"testing"

	"wordler"
	"wordler/wordlist"
)

func TestNew(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if s.Remaining() == 0 {
		t.Errorf("Only found %d words in initial list.", s.Remaining())
	}

	s, err = New(wordlist.KeepOnlyOption{Exp: regexp.MustCompile("^smile$")})
	if err != nil {
		t.Fatal(err)
	}
	if s.Remaining() != 1 {
		t.Errorf("Smile! Found %d words.", s.Remaining())
	}

	s, err = New(wordlist.DeleteOption{Exp: regexp.MustCompile("..")})
	if err != nil {
		t.Fatal(err)
	}
	if s.Remaining() != 26 {
		t.Errorf("Found %d 1-letter dictionary entries.", s.Remaining())
	}

	s, err = New(wordlist.DeleteOption{Exp: regexp.MustCompile("..")}, wordlist.KeepOnlyOption{Exp: regexp.MustCompile("..")})
	if err != nil {
		t.Fatal(err)
	}
	if s.Remaining() != 0 {
		t.Errorf("Found %d dictionary entries that are and are not 2 letter.", s.Remaining())
	}
}

func TestGuess(t *testing.T) {
	testList := []string{"foo", "bar", "bam", "zap", "zbz"}
	guesser := From(testList)
	guess := guesser.Guess()

	if !guesser.s.Contains(guess) {
		t.Errorf("Guess %v not found in solution list %#v", guess, guesser.s)
	}
	if !guesser.g.Contains(guess) {
		t.Errorf("Guess %v not found in guess list %#v", guess, guesser.g)
	}

	if guesser.s.Contains("bogus") {
		t.Errorf("How did \"bogus\" get in solution list %#v", guesser.s)
	}
	if guesser.g.Contains("bogus") {
		t.Errorf("How did \"bogus\" get in guess list %#v", guesser.g)
	}

	singleton := "foo"
	guesser.s = wordlist.New([]string{singleton})
	guesser.g = wordlist.New([]string{singleton})
	guess = guesser.Guess()
	if guess != singleton {
		t.Errorf("Want guess %v, got %v", singleton, guess)
	}

	guesser.s = nil
	guesser.g = nil
	guess = guesser.Guess()
	if guess != "" {
		t.Errorf("Want empty string, got %v", guess)
	}
}

func TestReact(t *testing.T) {
	testList := []string{"foo", "bar", "bam", "zap", "zbz"}
	cases := []struct {
		guess, response    string
		solutions, guesses *wordlist.WordList // valid solutions and guesses
	}{{
		guess:     "bar",
		response:  string([]byte{wordler.CORRECT, wordler.CORRECT, wordler.CORRECT}),
		solutions: wordlist.New([]string{"bar"}),
		guesses:   wordlist.New([]string{"bar"}),
	}, {
		guess:     "###",
		response:  string([]byte{wordler.NIL, wordler.NIL, wordler.NIL}),
		solutions: wordlist.New(testList),
		guesses:   wordlist.New(testList),
	}, {
		guess:     "abc",
		response:  string([]byte{wordler.NIL, wordler.NIL, wordler.NIL}),
		solutions: wordlist.New([]string{"foo"}),
		guesses:   wordlist.New(testList),
	}, {
		guess:     "b##",
		response:  string([]byte{wordler.ELSEWHERE, wordler.NIL, wordler.NIL}),
		solutions: wordlist.New([]string{"zbz"}),
		guesses:   wordlist.New([]string{"bar", "bam", "zbz"}),
	}, {
		guess:     "#oz",
		response:  string([]byte{wordler.NIL, wordler.NIL, wordler.ELSEWHERE}),
		solutions: wordlist.New([]string{"zap"}),
		guesses:   wordlist.New([]string{"zap", "zbz"}),
	}, {
		guess:     "zfo",
		response:  string([]byte{wordler.NIL, wordler.ELSEWHERE, wordler.CORRECT}),
		solutions: wordlist.New([]string{"foo"}),
		guesses:   wordlist.New([]string{"foo"}),
	}, {
		guess:     "b#r",
		response:  string([]byte{wordler.CORRECT, wordler.NIL, wordler.CORRECT}),
		solutions: wordlist.New([]string{"bar"}),
		guesses:   wordlist.New([]string{"bar"}),
	}}

	for _, c := range cases {
		t.Run(c.guess, func(t *testing.T) {
			guesser := From(testList)
			if err := guesser.React(c.guess, c.response); err != nil {
				t.Errorf("got error %v", err)
			}
			if want, got := c.solutions, guesser.s; !want.Equals(got) {
				t.Errorf("want %#v != got %#v", want, got)
			}
			if want, got := c.guesses, guesser.g; !want.Equals(got) {
				t.Errorf("want %#v != got %#v", want, got)
			}
		})
	}

	// guess and response are expected to be same length.
	// test response validation
	s := &Solver{}
	for _, r := range []string{"", "+", "__________", "++ _", "_*+ ", " **+"} {
		t.Run(r, func(t *testing.T) {
			if err := s.React("this", r); err == nil {
				t.Error("want error, got nil")
			}
		})
	}
}

func TestRemaining(t *testing.T) {
	s := &Solver{
		s: wordlist.New([]string{"f"}),
		g: wordlist.New([]string{"f", "g", "h"}),
	}
	if want, got := 1, s.Remaining(); want != got {
		t.Errorf("want %d, got %d", want, got)
	}

	s.s = wordlist.New([]string{})
	if want, got := 0, s.Remaining(); want != got {
		t.Errorf("want %d, got %d", want, got)
	}

	s.s = nil
	if want, got := 0, s.Remaining(); want != got {
		t.Errorf("want %d, got %d", want, got)
	}

	if want, got := 3, s.g.Length(); want != got {
		t.Errorf("want %d, got %d", want, got)
	}

	s = nil
	if want, got := 0, s.Remaining(); want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestDoubleLetters(t *testing.T) {
	// Test cases with double-letters in the guess
	cases := []struct {
		word, guess string
		response    []byte
	}{{
		word:     "forty",
		guess:    "worry",
		response: []byte{wordler.NIL, wordler.CORRECT, wordler.CORRECT, wordler.NIL, wordler.CORRECT},
	}, {
		word:     "forty",
		guess:    "robot",
		response: []byte{wordler.ELSEWHERE, wordler.CORRECT, wordler.NIL, wordler.NIL, wordler.ELSEWHERE},
	}, {
		word:     "foyer",
		guess:    "carer",
		response: []byte{wordler.NIL, wordler.NIL, wordler.NIL, wordler.CORRECT, wordler.CORRECT},
	}, {
		word:     "ab",
		guess:    "aa",
		response: []byte{wordler.CORRECT, wordler.NIL},
	}, {
		word:     "aab",
		guess:    "baa",
		response: []byte{wordler.ELSEWHERE, wordler.CORRECT, wordler.ELSEWHERE},
	}, {
		word:     "aab",
		guess:    "bab",
		response: []byte{wordler.NIL, wordler.CORRECT, wordler.CORRECT},
	}, {
		word:     "aab",
		guess:    "bba",
		response: []byte{wordler.ELSEWHERE, wordler.NIL, wordler.ELSEWHERE},
	}}

	for _, c := range cases {
		t.Run(c.guess, func(t *testing.T) {
			s := &Solver{s: wordlist.New([]string{c.word})}
			response := string(c.response)
			if err := s.React(c.guess, response); err != nil {
				t.Errorf("Error: %v", err)
			}
			if s.Remaining() != 1 {
				t.Errorf("want 1, got %d", s.Remaining())
			}
			want := wordlist.New([]string{c.word})
			if !want.Equals(s.s) {
				t.Errorf("want %#v; got %#v", want, s.s)
			}
		})
	}
}

func TestEliminateGuess(t *testing.T) {
	w := wordlist.New([]string{"array", "foray", "stray", "spray"})
	s := &Solver{
		s: w.Clone(),
		g: w.Clone(),
	}
	response := string([]byte{wordler.NIL, wordler.NIL, wordler.CORRECT, wordler.CORRECT, wordler.CORRECT})

	if err := s.React("stray", response); err != nil {
		t.Errorf("Error: %v", err)
	}
	if got, want := s.s, wordlist.New([]string{"array", "foray"}); !got.Equals(want) {
		t.Errorf("solutions: want %#v; got %#v", want, s.s)
	}
	if got, want := s.g, wordlist.New([]string{"array", "foray", "spray"}); !got.Equals(want) {
		t.Errorf("guesses: want %#v; got %#v", want, s.g)
	}

	if err := s.React("array", response); err != nil {
		t.Errorf("Error: %v", err)
	}
	if got, want := s.s, wordlist.New([]string{"foray"}); !got.Equals(want) {
		t.Errorf("solutions: want %#v; got %#v", want, s.s)
	}
	if got, want := s.g, wordlist.New([]string{"foray", "spray"}); !got.Equals(want) {
		t.Errorf("guesses: want %#v; got %#v", want, s.g)
	}
}

func TestNotInWordle(t *testing.T) {
	l := []string{"a", "b", "c"}
	s := &Solver{s: wordlist.New(l), g: wordlist.New(l)}

	s.NotInWordle("b")
	want := wordlist.New([]string{"a", "c"})

	if !want.Equals(s.s) {
		t.Errorf("solutions: want %#v; got %#v", want, s.s)
	}
	if !want.Equals(s.g) {
		t.Errorf("guesses: want %#v; got %#v", want, s.g)
	}

	s.NotInWordle("ac")
	if !want.Equals(s.s) {
		t.Errorf("solutions: want %#v; got %#v", want, s.s)
	}
	if !want.Equals(s.g) {
		t.Errorf("guesses: want %#v; got %#v", want, s.g)
	}
}
