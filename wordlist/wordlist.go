package wordlist

import (
	"bufio"
	"os"
	"reflect"
	"regexp"
	"sync"
	"unicode"
)

type WordList struct {
	words map[string]bool
}

// Loader is the DictionaryLoader; it is an exported variable so that tests can
// provide a platform-independent substitution.
var Loader DictionaryLoader = &platformLoader{}

// NewDictionary returns a *WordList containing /usr/share/dict/words.
func NewDictionary(options ...Option) (*WordList, error) {
	return Loader.Load(options...)
}

// DictionaryLoader is the interface for loading the dictionary into a WordList.
type DictionaryLoader interface {
	// Load loads the dictionary, filtered via the given options.
	Load(options ...Option) (*WordList, error)
}

// platformLoader loads dictionary on this platform.
type platformLoader struct {
	once           sync.Once
	fullDictionary []string
	err            error
}

// TODO: make this platform-independent via goos.Is*
func (p *platformLoader) Load(options ...Option) (*WordList, error) {
	p.once.Do(func() {
		var file *os.File
		file, p.err = os.Open(`/usr/share/dict/words`)
		if p.err != nil {
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
	SCAN:
		for scanner.Scan() {
			word := scanner.Text()
			for _, c := range word {
				if unicode.IsUpper(c) {
					continue SCAN
				}
			}
			p.fullDictionary = append(p.fullDictionary, word)
		}
		p.err = scanner.Err()
	})
	return New(p.fullDictionary, options...), p.err
}

// New creates a new WordList containing the words in s.
func New(s []string, options ...Option) *WordList {
	m := make(map[string]bool)
	for _, word := range s {
		m[word] = true
	}
	w := &WordList{m}
	for _, o := range options {
		o.apply(w)
	}
	return w
}

// Equals compares two WordLists and returns true if they contain the same
// words.
func (this *WordList) Equals(that *WordList) bool {
	switch {
	case this == nil:
		return that == nil
	case this.Length() != that.Length():
		return false
	case that == nil:
		return this == nil
	}

	return reflect.DeepEqual(this.words, that.words)
}

// Length returns the number of words in the list.
func (w *WordList) Length() int {
	if w == nil {
		return 0
	}
	return len(w.words)
}

// Delete removes all elements that match the given Regexp.
func (w *WordList) Delete(r *regexp.Regexp) {
	w.filter(r, true)
}

// KeepOnly removes all elements that don't match the given Regexp.
func (w *WordList) KeepOnly(r *regexp.Regexp) {
	w.filter(r, false)
}

// filter WordList based on r; if omit is true, delete matching items. If omit
// is false, keep matching items.
func (w *WordList) filter(r *regexp.Regexp, omit bool) {
	if w == nil {
		return
	}
	for word, _ := range w.words {
		if omit == r.MatchString(word) {
			delete(w.words, word)
		}
	}
}

// Contains returns true if word is in the WordList.
func (w *WordList) Contains(word string) bool {
	if w == nil {
		return false
	}
	return w.words[word]
}

// Random returns a random word from the WordList
func (w *WordList) Random() string {
	// map iteration goes in a random order!
	if w != nil {
		for guess, _ := range w.words {
			return guess
		}
	}
	return ""
}

// Option represents a constraint to place on a WordList.
type Option interface {
	apply(*WordList)
}

// KeepOnlyOption specifies that Solver should only include words that match
// the given expression.
type KeepOnlyOption struct {
	Exp *regexp.Regexp
}

// apply fulfills the Option interface
func (k KeepOnlyOption) apply(w *WordList) {
	w.KeepOnly(k.Exp)
}

// DeleteOption specifies that Solver should exclude words that match
// the given expression.
type DeleteOption struct {
	Exp *regexp.Regexp
}

// apply fulfills the Option interface
func (d DeleteOption) apply(w *WordList) {
	w.Delete(d.Exp)
}
