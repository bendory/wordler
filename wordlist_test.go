package wordler

import (
	"regexp"
	"testing"
)

var baseList = []string{"foo", "bar", "bam", "zoo"}

func TestLength(t *testing.T) {
	cases := []struct {
		w *WordList
		l int
	}{
		{&WordList{}, 0},
		{nil, 0},
		{NewWordList(baseList), len(baseList)},
		{NewWordList(baseList), len(baseList)},
		{&WordList{nil}, 0},
		{NewWordList([]string{}), 0},
	}

	for _, c := range cases {
		if got, want := c.w.Length(), c.l; got != want {
			t.Errorf("%#v: got %d, want %d", c.w, got, want)
		}
	}
}

func TestEquals(t *testing.T) {
	cases := []struct {
		left, right *WordList
		want        bool
	}{
		{nil, nil, true},
		{nil, &WordList{}, false},
		{nil, NewWordList(baseList), false},
		{NewWordList(baseList), NewWordList(baseList), true},
		{NewWordList(baseList), NewWordList([]string{"foo"}), false},
	}

	for _, c := range cases {
		if got, want := c.left.Equals(c.right), c.want; got != want {
			t.Errorf("got %t, want %t: %#v <=> %#v", got, want, c.left, c.right)
		}
		if got, want := c.right.Equals(c.left), c.want; got != want {
			t.Errorf("got %t, want %t: %#v <=> %#v", got, want, c.right, c.left)
		}
	}
}

func TestDelete(t *testing.T) {
	empty := NewWordList([]string{})
	cases := []struct {
		filter string
		want   *WordList
	}{
		{".oo", NewWordList([]string{"bar", "bam"})},
		{"ba", NewWordList([]string{"foo", "zoo"})},
		{"....", NewWordList(baseList)},
		{"...", empty},
		{"..", empty},
		{"", empty},
		{"nomatch", NewWordList(baseList)},
	}

	for _, c := range cases {
		t.Run(c.filter, func(t *testing.T) {
			w := NewWordList(baseList)
			w.Delete(regexp.MustCompile(c.filter))

			if got, want := w.Length(), c.want.Length(); got != want {
				t.Fatalf("got %d words %#v, want %d %#v", got, w, want, c.want)
			}

			if !w.Equals(c.want) {
				t.Errorf("got %#v, want %#v", w, c.want)
			}
		})
	}
}

func TestKeepOnly(t *testing.T) {
	empty := NewWordList([]string{})
	cases := []struct {
		filter string
		want   *WordList
	}{
		{".oo", NewWordList([]string{"foo", "zoo"})},
		{"ba", NewWordList([]string{"bar", "bam"})},
		{"....", empty},
		{"...", NewWordList(baseList)},
		{"..", NewWordList(baseList)},
		{"", NewWordList(baseList)},
		{"nomatch", empty},
	}

	for _, c := range cases {
		t.Run(c.filter, func(t *testing.T) {
			w := NewWordList(baseList)
			w.KeepOnly(regexp.MustCompile(c.filter))

			if got, want := w.Length(), c.want.Length(); got != want {
				t.Fatalf("got %d words %#v, want %d %#v", got, w, want, c.want)
			}

			if !w.Equals(c.want) {
				t.Errorf("got %#v, want %#v", w, c.want)
			}
		})
	}
}

func TestNewWordList(t *testing.T) {
	l := baseList[:]

	w := NewWordList(l)
	want := NewWordList(baseList)
	if !w.Equals(want) {
		t.Errorf("got %#v, want %#v", w, want)
	}

	// Changing w should have no effect on l.
	w2 := NewWordList(l)
	w.Delete(regexp.MustCompile("."))
	if w.Equals(w2) {
		t.Errorf("%#v should not equal %#v", w, w2)
	}

	// Changing l should have no effect on w.
	w = NewWordList(l)
	l[0] = l[0] + " bogus"
	if !w.Equals(want) {
		t.Errorf("got %#v, want %#v", w, want)
	}
	if !w.Equals(w2) {
		t.Errorf("got %#v, want %#v", w, w2)
	}
}
