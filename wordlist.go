package wordler

import (
	"reflect"
	"regexp"
)

type WordList struct {
	words map[string]bool
}

// NewWordList creates a new WordList containing the words in s.
func NewWordList(s []string) *WordList {
	m := make(map[string]bool)
	for _, word := range s {
		m[word] = true
	}
	return &WordList{m}
}

// Equals compares two WordLists and returns true if they contain the same
// words. Note that comparison is independent of order.
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
