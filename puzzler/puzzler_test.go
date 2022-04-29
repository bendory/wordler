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
	was := wordlist.Loader
	defer func() { wordlist.Loader = was }()
	wordlist.Loader = fakeLoader{words: list}

	p, err := New(true)
	if err != nil {
		t.Errorf("error: %v", err)
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
		list        []string
		guess, word string
		response    *Response
		err         error
	}{{
		list:     []string{"foo", "bar", "bam", "zap"},
		guess:    "bar",
		word:     "foo",
		response: &Response{Detail: string([]byte{wordler.NIL, wordler.NIL, wordler.NIL}), WordsRemaining: 1},
	}, {
		list:     []string{"foo", "bar", "bam", "zap"},
		guess:    "foo",
		word:     "foo",
		response: &Response{Detail: string([]byte{wordler.CORRECT, wordler.CORRECT, wordler.CORRECT}), WordsRemaining: 1},
	}, {
		list:     []string{"foo", "bar", "bam", "zap"},
		guess:    "bam",
		word:     "bar",
		response: &Response{Detail: string([]byte{wordler.CORRECT, wordler.CORRECT, wordler.NIL}), WordsRemaining: 1},
	}, {
		list:     []string{"foo", "bar", "bam", "zap", "pta"},
		guess:    "pta",
		word:     "zap",
		response: &Response{Detail: string([]byte{wordler.ELSEWHERE, wordler.NIL, wordler.ELSEWHERE}), WordsRemaining: 1},
	}}

	for _, c := range cases {
		t.Run(fmt.Sprint("g:", c.guess, "+w:", c.word), func(t *testing.T) {
			was := wordlist.Loader
			defer func() { wordlist.Loader = was }()
			wordlist.Loader = fakeLoader{words: c.list}

			p, err := New(true)
			if err != nil {
				t.Errorf("error: %v", err)
			}
			p.word = c.word
			response, err := p.Guess(c.guess)
			if err != c.err {
				t.Errorf("want %v; got %v", c.err, err)
			}
			if want, got := c.response.Detail, response.Detail; want != got {
				t.Errorf("want %v; got %v", want, got)
			}
			if want, got := c.response.WordsRemaining, response.WordsRemaining; want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		})
	}
}
