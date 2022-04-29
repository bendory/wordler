package wordler

import (
	"errors"
	"testing"

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
	list := []string{"foo", "bar", "bam", "zap"}
	was := wordlist.Loader
	defer func() { wordlist.Loader = was }()
	wordlist.Loader = fakeLoader{words: list}

	p, err := New(true)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	p.word = list[0]

	gotResponse, gotErr := p.Guess(list[1])
	if gotErr != nil {
		t.Errorf("error: %v", err)
	}
	if want, got := "___", gotResponse.Detail; want != got {
		t.Errorf("want %#v; got %#v", want, got)
	}
	// TODO: check other values in gotResponse

	gotResponse, gotErr = p.Guess(list[0])
	if gotErr != nil {
		t.Errorf("error: %v", err)
	}
	if want, got := "+++", gotResponse.Detail; want != got {
		t.Errorf("want %#v; got %#v", want, got)
	}
	// TODO: check other values in gotResponse

	p.word = list[1]
	gotResponse, gotErr = p.Guess(list[2])
	if gotErr != nil {
		t.Errorf("error: %v", err)
	}
	if want, got := "++_", gotResponse.Detail; want != got {
		t.Errorf("want %#v; got %#v", want, got)
	}
	// TODO: check other values in gotResponse

	// TODO: test ELSEWHERE values too
}
