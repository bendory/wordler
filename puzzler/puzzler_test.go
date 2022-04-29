package puzzler

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

	p, err := New(&Args{Hard: true, WordLength: 3})
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
	p, err = New(nil)
	if err == nil {
		t.Error("want error, got nil")
	}

	f.err = nil
	p, err = New(&Args{WordLength: 0})
	if want, got := NoWordsRemainingErr, err; want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestValidate(t *testing.T) {
	list := []string{"foo", "bar", "bam", "zap"}
	p := Wordle{
		dict:             wordlist.New(list),
		remaining:        wordlist.New(list),
		word:             list[0],
		remainingGuesses: wordler.DEFAULT_GUESSES,
		hard:             true,
	}
	if got := p.validate(list[0]); got != nil {
		t.Errorf("want nil; got %#v", got)
	}
	if want, got := NotInDictionaryErr, p.validate("bogus"); want != got {
		t.Errorf("want %#v; got %#v", want, got)
	}

	p.remaining = wordlist.New([]string{""})
	if want, got := InvalidGuessErr, p.validate(list[0]); want != got {
		t.Errorf("want %#v; got %#v", want, got)
	}
	p.hard = false
	if got := p.validate(list[0]); got != nil {
		t.Errorf("want nil; got %#v", got)
	}

	p.remaining = nil
	if want, got := NoWordsRemainingErr, p.validate(list[0]); want != got {
		t.Errorf("want %#v; got %#v", want, got)
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
				remainingGuesses: wordler.DEFAULT_GUESSES,
				hard:             true,
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

func TestRemaining(t *testing.T) {
	list := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	p := Wordle{
		dict:             wordlist.New(list),
		remaining:        wordlist.New(list),
		word:             list[len(list)-1],
		remainingGuesses: wordler.DEFAULT_GUESSES,
		hard:             true,
	}

	if p.Guesses() >= p.Words() {
		t.Errorf("want guesses %d < remaining words %d", p.Guesses(), p.Words())
	}

	// A rejected guess should not consume a remaining guess.
	response, err := p.Guess("not in dictionary")
	if response != "" {
		t.Errorf("want empty string, got %v", response)
	}
	if err != NotInDictionaryErr {
		t.Errorf("want %v, got %v", NotInDictionaryErr, err)
	}

	legit := list[:wordler.DEFAULT_GUESSES]
	bogus := list[wordler.DEFAULT_GUESSES:]
	if len(legit) < 1 {
		t.Error("legit is bogus! want at least 1 item.")
	}
	if len(bogus) < 1 {
		t.Error("bogus is bogus! want at least 1 item.")
	}

	for i, guess := range legit {
		t.Run(fmt.Sprintf("legit_%d", i), func(t *testing.T) {
			if want, got := wordler.DEFAULT_GUESSES-i, p.Guesses(); want != got {
				t.Errorf("pre-guess: want %d guesses, got %d", want, got)
			}
			response, err = p.Guess(guess)
			if want, got := string(wordler.NIL), response; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			if err != nil {
				t.Errorf("want nil, got %v", err)
			}
			if want, got := len(list)-i-1, p.Words(); want != got {
				t.Errorf("want %d remaining words, got %d", want, got)
			}
			if want, got := wordler.DEFAULT_GUESSES-i-1, p.Guesses(); want != got {
				t.Errorf("post-guess: want %d guesses, got %d", want, got)
			}
		})
	}

	if want, got := 0, p.Guesses(); want != got {
		t.Errorf("want %d guesses, got %d", want, got)
	}
	if want, got := len(bogus), p.Words(); want != got {
		t.Errorf("want %d remaining words, got %d", want, got)
	}

	for i, guess := range bogus {
		t.Run(fmt.Sprintf("bogus_%d", i), func(t *testing.T) {
			if want, got := 0, p.Guesses(); want != got {
				t.Errorf("pre-guess: want %d guesses, got %d", want, got)
			}
			response, err = p.Guess(guess)
			if want, got := "", response; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			if err != OutOfGuessesErr {
				t.Errorf("want nil, got %v", err)
			}
			if want, got := len(list)-len(legit), p.Words(); want != got {
				t.Errorf("want %d remaining words, got %d", want, got)
			}
			if want, got := 0, p.Guesses(); want != got {
				t.Errorf("post-guess: want %d guesses, got %d", want, got)
			}
		})
	}

	response, err = p.Guess("another bogus guess")
	if response != "" {
		t.Errorf("got response %v to bogus guess", response)
	}
	if want, got := OutOfGuessesErr, err; want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestGiveUp(t *testing.T) {
	list := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	p := Wordle{
		dict:             wordlist.New(list),
		remaining:        wordlist.New(list),
		word:             list[0],
		remainingGuesses: wordler.DEFAULT_GUESSES,
	}

	if want, got := wordler.DEFAULT_GUESSES, p.Guesses(); want != got {
		t.Errorf("want %v, got %v", want, got)
	}
	if want, got := len(list), p.Words(); want != got {
		t.Errorf("want %v, got %v", want, got)
	}
	response, err := p.Guess(list[1])
	if want, got := string(wordler.NIL), response; want != got {
		t.Errorf("want %v, got %v", want, got)
	}
	if err != nil {
		t.Errorf("want nil, got %v", err)
	}

	// Now GiveUp().
	if want, got := p.word, p.GiveUp(); want != got {
		t.Errorf("want %v, got %v", want, got)
	}
	if want, got := 0, p.Guesses(); want != got {
		t.Errorf("want %v, got %v", want, got)
	}
	if want, got := 1, p.Words(); want != got {
		t.Errorf("want %v, got %v", want, got)
	}

	response, err = p.Guess(list[1])
	if response != "" {
		t.Errorf("want empty string, got %v", response)
	}
	if want, got := OutOfGuessesErr, err; want != got {
		t.Errorf("want %v, got %v", want, got)
	}
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

func TestSolution(t *testing.T) {
	list := []string{"foo", "bar", "bam", "zap"}
	was := wordlist.Loader
	defer func() { wordlist.Loader = was }()
	f := &fakeLoader{words: list}
	wordlist.Loader = f

	args := &Args{Solution: list[0], WordLength: 3}
	if _, err := New(args); err != nil {
		t.Errorf("error: %v", err)
	}
	args.Solution = "bogus"
	if _, err := New(args); err != NotInDictionaryErr {
		t.Errorf("error: %v", err)
	}
}

func TestElsewhere(t *testing.T) {
	// This test case came from a bug found in an actual run.
	guess, word := "dreamt", "machin"
	list := []string{guess, word, "ambush"}
	p := Wordle{
		dict:             wordlist.New(list),
		remaining:        wordlist.New(list),
		word:             word,
		remainingGuesses: wordler.DEFAULT_GUESSES,
	}

	response, err := p.Guess(guess)
	if err != nil {
		t.Errorf("want nil, got %v", err)
	}
	if want, got := string([]byte{wordler.NIL, wordler.NIL, wordler.NIL, wordler.ELSEWHERE, wordler.ELSEWHERE, wordler.NIL}), response; want != got {
		t.Errorf("want %v, got %v", want, got)
	}
	if want, got := 2, p.Words(); want != got {
		t.Errorf("want %d words remaining, got %d", want, got)
	}
}

func TestDoubleLetters(t *testing.T) {
	// This test case came from a bug found in an actual run.
	guesses := []string{"huffle", "whimmy"}
	word := "machin"

	cases := []struct {
		guess, response string
		remaining       int
	}{{
		guess:     guesses[0],
		response:  string([]byte{wordler.ELSEWHERE, wordler.NIL, wordler.NIL, wordler.NIL, wordler.NIL, wordler.NIL}),
		remaining: 2,
	}, {
		guess:     guesses[1],
		response:  string([]byte{wordler.NIL, wordler.ELSEWHERE, wordler.ELSEWHERE, wordler.ELSEWHERE, wordler.NIL, wordler.NIL}),
		remaining: 1,
	}}
	list := append(guesses, word)

	p := Wordle{
		dict:             wordlist.New(list),
		remaining:        wordlist.New(list),
		word:             word,
		remainingGuesses: wordler.DEFAULT_GUESSES,
	}

	for _, c := range cases {
		t.Run(c.guess, func(t *testing.T) {
			response, err := p.Guess(c.guess)
			if err != nil {
				t.Errorf("want nil, got %v", err)
			}
			if want, got := c.response, response; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			if want, got := c.remaining, p.Words(); want != got {
				t.Errorf("want %d words remaining, got %d", want, got)
			}
		})
	}
}
