package solver

import (
	"regexp"
	"strings"
	"testing"

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

	s, err = New(KeepOnlyOption{regexp.MustCompile("^smile$")})
	if err != nil {
		t.Fatal(err)
	}
	if s.Remaining() != 1 {
		t.Errorf("Smile! Found %d words.", s.Remaining())
	}

	s, err = New(DeleteOption{regexp.MustCompile("..")})
	if err != nil {
		t.Fatal(err)
	}
	if s.Remaining() != 26 {
		t.Errorf("Found %d 1-letter dictionary entries.", s.Remaining())
	}

	s, err = New(DeleteOption{regexp.MustCompile("..")}, KeepOnlyOption{regexp.MustCompile("..")})
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
		{"bar", strings.Repeat(string(CORRECT), 3), wordlist.New([]string{"bar"})},
		{"###", strings.Repeat(string(NIL), 3), wordlist.New(testList)},
		{"abc", strings.Repeat(string(NIL), 3), wordlist.New([]string{"foo"})},
		{
			"b##",
			string(ELSEWHERE) + string(NIL) + string(NIL),
			wordlist.New([]string{"zbz"}),
		}, {
			"#oz",
			string(NIL) + string(NIL) + string(ELSEWHERE),
			wordlist.New([]string{"zap"}),
		}, {
			"zfo",
			string(NIL) + string(ELSEWHERE) + string(CORRECT),
			wordlist.New([]string{"foo"}),
		}, {
			"b#r",
			string(CORRECT) + string(NIL) + string(CORRECT),
			wordlist.New([]string{"bar"}),
		},
	}

	for _, c := range cases {
		t.Run(c.guess, func(t *testing.T) {
			guesser := &Solver{w: wordlist.New(testList)}
			guesser.React(c.guess, c.response)
			if want, got := c.remaining, guesser.w; !want.Equals(got) {
				t.Errorf("want %#v != got %#v", want, got)
			}
		})
	}
}

func TestRemaining(t *testing.T) {
	g := &Solver{w: wordlist.New([]string{"f"})}
	if want, got := 1, g.Remaining(); want != got {
		t.Errorf("want %d, got %d", want, got)
	}
	g.w = wordlist.New([]string{})
	if want, got := 0, g.Remaining(); want != got {
		t.Errorf("want %d, got %d", want, got)
	}
	g.w = nil
	if want, got := 0, g.Remaining(); want != got {
		t.Errorf("want %d, got %d", want, got)
	}
	g = nil
	if want, got := 0, g.Remaining(); want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestDoubleLetters(t *testing.T) {
	// Test cases with double-letters in the guess
	cases := []struct {
		guess    string
		response []byte
	}{
		{
			guess:    "worry",
			response: []byte{NIL, CORRECT, CORRECT, NIL, CORRECT},
		}, {
			guess:    "robot",
			response: []byte{ELSEWHERE, CORRECT, NIL, NIL, ELSEWHERE},
		},
	}

	for _, c := range cases {
		t.Run(c.guess, func(t *testing.T) {
			g := &Solver{w: wordlist.New([]string{"forty"})}
			response := string(c.response)
			g.React(c.guess, response)
			if g.Remaining() != 1 {
				t.Errorf("want 1, got %d", g.Remaining())
			}
			want := wordlist.New([]string{"forty"})
			if !want.Equals(g.w) {
				t.Errorf("want %#v; got %#v", want, g.w)
			}
		})
	}

	// Guess case should be eliminated.
	g := &Solver{
		w: wordlist.New([]string{"array", "foray"}),
		have: map[byte]bool{
			'r': true,
			'a': true,
			'y': true,
		},
	}
	response := string([]byte{NIL, NIL, CORRECT, CORRECT, CORRECT})
	g.React("array", response)
	want := wordlist.New([]string{"foray"})
	if !want.Equals(g.w) {
		t.Errorf("want %#v; got %#v", want, g.w)
	}
}
