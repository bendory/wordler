package guesser

import (
	"regexp"
	"strings"
	"testing"

	"wordler/wordlist"
)

func TestNew(t *testing.T) {
	g, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if g.Remaining() == 0 {
		t.Errorf("Only found %d words in initial list.", g.Remaining())
	}

	g, err = New(KeepOnlyOption{regexp.MustCompile("^smile$")})
	if err != nil {
		t.Fatal(err)
	}
	if g.Remaining() != 1 {
		t.Errorf("Smile! Found %d words.", g.Remaining())
	}

	g, err = New(DeleteOption{regexp.MustCompile("..")})
	if err != nil {
		t.Fatal(err)
	}
	if g.Remaining() != 26 {
		t.Errorf("Found %d 1-letter dictionary entries.", g.Remaining())
	}

	g, err = New(DeleteOption{regexp.MustCompile("..")}, KeepOnlyOption{regexp.MustCompile("..")})
	if err != nil {
		t.Fatal(err)
	}
	if g.Remaining() != 0 {
		t.Errorf("Found %d dictionary entries that are and are not 2 letter.", g.Remaining())
	}
}

func TestGuess(t *testing.T) {
	testList := []string{"foo", "bar", "bam", "zap", "zbz"}
	guesser := &Guesser{w: wordlist.New(testList)}
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
		{"bar", strings.Repeat(string(RIGHT_LETTER_RIGHT_PLACE), 3), wordlist.New([]string{"bar"})},
		{"###", strings.Repeat(string(LETTER_NOT_IN_WORD), 3), wordlist.New(testList)},
		{"abc", strings.Repeat(string(LETTER_NOT_IN_WORD), 3), wordlist.New([]string{"foo"})},
		{
			"b##",
			string(RIGHT_LETTER_WRONG_PLACE) + string(LETTER_NOT_IN_WORD) + string(LETTER_NOT_IN_WORD),
			wordlist.New([]string{"zbz"}),
		}, {
			"#oz",
			string(LETTER_NOT_IN_WORD) + string(LETTER_NOT_IN_WORD) + string(RIGHT_LETTER_WRONG_PLACE),
			wordlist.New([]string{"zap"}),
		}, {
			"zfo",
			string(LETTER_NOT_IN_WORD) + string(RIGHT_LETTER_WRONG_PLACE) + string(RIGHT_LETTER_RIGHT_PLACE),
			wordlist.New([]string{"foo"}),
		}, {
			"b#r",
			string(RIGHT_LETTER_RIGHT_PLACE) + string(LETTER_NOT_IN_WORD) + string(RIGHT_LETTER_RIGHT_PLACE),
			wordlist.New([]string{"bar"}),
		},
	}

	for _, c := range cases {
		t.Run(c.guess, func(t *testing.T) {
			guesser := &Guesser{w: wordlist.New(testList)}
			guesser.React(c.guess, c.response)
			if want, got := c.remaining, guesser.w; !want.Equals(got) {
				t.Errorf("want %#v != got %#v", want, got)
			}
		})
	}
}

func TestRemaining(t *testing.T) {
	g := &Guesser{w: wordlist.New([]string{"f"})}
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