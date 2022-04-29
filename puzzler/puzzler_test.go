package wordler

import (
	"errors"
	"fmt"
	"testing"

	"wordler"
	"wordler/wordlist"
)

type fakeLoader struct {
	words []string // words to load
	err   error    // optional error for testing error paths
}

func (f fakeLoader) Load(options ...wordlist.Option) (*wordlist.WordList, error) {
	return wordlist.New(f.words, options...), f.err
}

func TestNew(t *testing.T) {
	list := []string{"foo", "bar", "bam", "zap"}
	was := wordlist.Loader
	defer func() { wordlist.Loader = was }()
	f := &fakeLoader{words: list}
	wordlist.Loader = f

	p, err := New(true)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if want, got := len(list), p.dict.Length(); want != got {
		t.Errorf("want %#v; got %#v", list, p.dict)
	}
	if !p.dict.Equals(p.remaining) {
		t.Errorf("%#v != %#v", p.dict, p.remaining)
	}
	if !p.dict.Contains(p.word) {
		t.Errorf("dict is missing %v", p.word)
	}

	f.err = errors.New("some error")
	p, err = New(true)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestValidate(t *testing.T) {
	list := []string{"foo", "bar", "bam", "zap"}
	p := Wordle{
		dict:             wordlist.New(list),
		remaining:        wordlist.New(list),
		word:             list[0],
		remainingGuesses: defaultGuesses,
		strict:           true,
	}
	if got := p.validate(list[0]); got != nil {
		t.Errorf("want nil; got %#v", got)
	}
	if want, got := NotInDictionaryErr, p.validate("bogus"); want != got {
		t.Errorf("want %#v; got %#v", want, got)
	}
	p.remaining = wordlist.New([]string{})
	if want, got := InvalidGuessErr, p.validate(list[0]); want != got {
		t.Errorf("want %#v; got %#v", want, got)
	}

	p.strict = false
	if got := p.validate(list[0]); got != nil {
		t.Errorf("want nil; got %#v", got)
	}
}

func TestGuess(t *testing.T) {
	cases := []struct {
		list                  []string
		guess, word, response string
		wordsRemaining        int
		err                   error
	}{{
		list:           []string{"foo", "bar", "bam", "zap"},
		guess:          "bar",
		word:           "foo",
		response:       string([]byte{wordler.NIL, wordler.NIL, wordler.NIL}),
		wordsRemaining: 1,
	}, {
		list:           []string{"foo", "bar", "bam", "zap"},
		guess:          "foo",
		word:           "foo",
		response:       string([]byte{wordler.CORRECT, wordler.CORRECT, wordler.CORRECT}),
		wordsRemaining: 1,
	}, {
		list:           []string{"foo", "bar", "bam", "zap"},
		guess:          "bam",
		word:           "bar",
		response:       string([]byte{wordler.CORRECT, wordler.CORRECT, wordler.NIL}),
		wordsRemaining: 1,
	}, {
		list:           []string{"foo", "bar", "bam", "zap", "pta"},
		guess:          "pta",
		word:           "zap",
		response:       string([]byte{wordler.ELSEWHERE, wordler.NIL, wordler.ELSEWHERE}),
		wordsRemaining: 1,
	}, {
		list:  []string{"a"},
		guess: "b",
		word:  "a",
		err:   NotInDictionaryErr,
	}}

	for _, c := range cases {
		t.Run(fmt.Sprint("g:", c.guess, "+w:", c.word), func(t *testing.T) {
			p := Wordle{
				dict:             wordlist.New(c.list),
				remaining:        wordlist.New(c.list),
				word:             c.word,
				remainingGuesses: defaultGuesses,
				strict:           true,
			}

			if want, got := len(c.list), p.Words(); want != got {
				t.Errorf("want %d words, got %d", want, got)
			}

			response, err := p.Guess(c.guess)
			if err != c.err {
				t.Errorf("want %v; got %v", c.err, err)
			} else if err != nil {
				return
			}
			if want, got := c.response, response; want != got {
				t.Errorf("want %v; got %v", want, got)
			}
			if want, got := c.wordsRemaining, p.Words(); want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		})
	}
}

func TestRemainingGuesses(t *testing.T) {
	// TODO: add tests
}

func TestGiveUp(t *testing.T) {
	// TODO: add tests
}

func TestNilWordle(t *testing.T) {
	var w *Wordle
	r, err := w.Guess("foo")
	if r != "" {
		t.Errorf("want empty string, got %v", r)
	}
	if err != OutOfGuessesErr {
		t.Errorf("want %v, got %v", OutOfGuessesErr, err)
	}
	if want, got := 0, w.Guesses(); want != got {
		t.Errorf("want %v, got %v", want, got)
	}
	if want, got := "", w.GiveUp(); want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}
