package wordlist

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
		{New(baseList), len(baseList)},
		{New(baseList), len(baseList)},
		{&WordList{nil}, 0},
		{New([]string{}), 0},
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
		{nil, New(baseList), false},
		{New(baseList), New(baseList), true},
		{New(baseList), New([]string{"foo"}), false},
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
	empty := New([]string{})
	cases := []struct {
		filter string
		want   *WordList
	}{
		{".oo", New([]string{"bar", "bam"})},
		{"ba", New([]string{"foo", "zoo"})},
		{"....", New(baseList)},
		{"...", empty},
		{"..", empty},
		{"", empty},
		{"nomatch", New(baseList)},
	}

	for _, c := range cases {
		t.Run(c.filter, func(t *testing.T) {
			w := New(baseList)
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
	empty := New([]string{})
	cases := []struct {
		filter string
		want   *WordList
	}{
		{".oo", New([]string{"foo", "zoo"})},
		{"ba", New([]string{"bar", "bam"})},
		{"....", empty},
		{"...", New(baseList)},
		{"..", New(baseList)},
		{"", New(baseList)},
		{"nomatch", empty},
	}

	for _, c := range cases {
		t.Run(c.filter, func(t *testing.T) {
			w := New(baseList)
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

func TestNew(t *testing.T) {
	l := baseList[:]

	w := New(l)
	want := New(baseList)
	if !w.Equals(want) {
		t.Errorf("got %#v, want %#v", w, want)
	}

	// Changing w should have no effect on l.
	w2 := New(l)
	w.Delete(regexp.MustCompile("."))
	if w.Equals(w2) {
		t.Errorf("%#v should not equal %#v", w, w2)
	}

	// Changing l should have no effect on w.
	w = New(l)
	l[0] = l[0] + " bogus"
	if !w.Equals(want) {
		t.Errorf("got %#v, want %#v", w, want)
	}
	if !w.Equals(w2) {
		t.Errorf("got %#v, want %#v", w, w2)
	}
}

func TestDictionary(t *testing.T) {
	d, err := NewDictionary()
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}

	if d.Length() < 10000 {
		t.Errorf("Dictionary looks small: only found %d words.", d.Length())
	}
}

func TestContains(t *testing.T) {
	d, err := NewDictionary()
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}

	if !d.Contains("toner") {
		t.Errorf("Dictionary does not contain toner.")
	}
	if d.Contains("not a word") {
		t.Errorf("Dictionary contains 'not a word'.")
	}
}
