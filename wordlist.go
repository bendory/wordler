package wordler

import (
	"reflect"
	"regexp"
)

type WordList struct {
	words []string
}

// New creates a new WordList containing the words in s.
func New(s []string) *WordList {
	return &WordList{words: s[:]} // copy s
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

	thisM := make(map[string]bool)
	thatM := make(map[string]bool)

	for _, w := range this.words {
		thisM[w] = true
	}
	for _, w := range that.words {
		thatM[w] = true
	}

	return reflect.DeepEqual(thisM, thatM)
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
	if w == nil {
		return
	}
	last := w.Length() - 1
	for i := last; i >= 0; i-- {
		if r.MatchString(w.words[i]) {
			w.words[i], w.words[last] = w.words[last], w.words[i]
			last--
		}
	}
	w.words = w.words[:last+1]
}
