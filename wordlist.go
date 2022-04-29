package wordler

import "regexp"

type WordList []string

// Filter filters all elements that match the given Regexp, returning a new
// WordList with whatever is left. Note that the new WordList is a slice of the
// original, so future modifications to the original will apply to corresponding
// entries in the returned WordList!
func (w WordList) Filter(r *regexp.Regexp) WordList {
	last := len(w) - 1
	for i := last; i >= 0; i-- {
		if r.MatchString(w[i]) {
			w[i], w[last] = w[last], w[i]
			last--
		}
	}
	return w[:last+1]
}
