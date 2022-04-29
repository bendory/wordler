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
	guesser := &Solver{w: wordlist.New(testList)}
	guess := guesser.Guess()

	if !guesser.w.Contains(guess) {
		t.Errorf("Guess %v not found in wordlist %#v", guess, guesser.w)
	}
	if guesser.w.Contains("bogus") {
		t.Errorf("How did \"bogus\" get in wordlist %#v", guesser.w)
	}

	singleton := "foo"
	guesser.w = wordlist.New([]string{singleton})
	guess = guesser.Guess()
	if guess != singleton {
		t.Errorf("Want guess %v, got %v", singleton, guess)
	}

	guesser.w = nil
	guess = guesser.Guess()
	if guess != "" {
		t.Errorf("Want empty string, got %v", guess)
	}
}

func TestReact(t *testing.T) {
	testList := []string{"foo", "bar", "bam", "zap", "zbz"}
	cases := []struct {
		guess, response string
		remaining       *wordlist.WordList
	}{
		{"bar", string([]byte{wordler.CORRECT, wordler.CORRECT, wordler.CORRECT}), wordlist.New([]string{"bar"})},
		{"###", string([]byte{wordler.NIL, wordler.NIL, wordler.NIL}), wordlist.New(testList)},
		{"abc", string([]byte{wordler.NIL, wordler.NIL, wordler.NIL}), wordlist.New([]string{"foo"})},
		{
			"b##",
			string([]byte{wordler.ELSEWHERE, wordler.NIL, wordler.NIL}),
			wordlist.New([]string{"zbz"}),
		}, {
			"#oz",
			string([]byte{wordler.NIL, wordler.NIL, wordler.ELSEWHERE}),
			wordlist.New([]string{"zap"}),
		}, {
			"zfo",
			string([]byte{wordler.NIL, wordler.ELSEWHERE, wordler.CORRECT}),
			wordlist.New([]string{"foo"}),
		}, {
			"b#r",
			string([]byte{wordler.CORRECT, wordler.NIL, wordler.CORRECT}),
			wordlist.New([]string{"bar"}),
		},
	}

	for _, c := range cases {
		t.Run(c.guess, func(t *testing.T) {
			guesser := &Solver{w: wordlist.New(testList)}
			if err := guesser.React(c.guess, c.response); err != nil {
				t.Errorf("got error %v", err)
			}
			if want, got := c.remaining, guesser.w; !want.Equals(got) {
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
	s := &Solver{w: wordlist.New([]string{"f"})}
	if want, got := 1, s.Remaining(); want != got {
		t.Errorf("want %d, got %d", want, got)
	}
	s.w = wordlist.New([]string{})
	if want, got := 0, s.Remaining(); want != got {
		t.Errorf("want %d, got %d", want, got)
	}
	s.w = nil
	if want, got := 0, s.Remaining(); want != got {
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
	}}

	for _, c := range cases {
		t.Run(c.guess, func(t *testing.T) {
			s := &Solver{w: wordlist.New([]string{c.word})}
			response := string(c.response)
			if err := s.React(c.guess, response); err != nil {
				t.Errorf("Error: %v", err)
			}
			if s.Remaining() != 1 {
				t.Errorf("want 1, got %d", s.Remaining())
			}
			want := wordlist.New([]string{c.word})
			if !want.Equals(s.w) {
				t.Errorf("want %#v; got %#v", want, s.w)
			}
		})
	}

	// Guess case should be eliminated.
	s := &Solver{
		w: wordlist.New([]string{"array", "foray"}),
		have: map[byte]bool{
			'r': true,
			'a': true,
			'y': true,
		},
	}
	response := string([]byte{wordler.NIL, wordler.NIL, wordler.CORRECT, wordler.CORRECT, wordler.CORRECT})
	if err := s.React("array", response); err != nil {
		t.Errorf("Error: %v", err)
	}
	want := wordlist.New([]string{"foray"})
	if !want.Equals(s.w) {
		t.Errorf("want %#v; got %#v", want, s.w)
	}
}

func TestNotInWordle(t *testing.T) {
	s := &Solver{w: wordlist.New([]string{"a", "b", "c"})}

	s.NotInWordle("b")
	want := wordlist.New([]string{"a", "c"})

	if !want.Equals(s.w) {
		t.Errorf("want %#v; got %#v", want, s.w)
	}

	s.NotInWordle("ac")
	if !want.Equals(s.w) {
		t.Errorf("want %#v; got %#v", want, s.w)
	}
}
