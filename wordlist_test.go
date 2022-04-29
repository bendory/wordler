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
		{&WordList{baseList}, len(baseList)},
		{&WordList{baseList}, len(baseList)},
		{&WordList{nil}, 0},
		{&WordList{[]string{}}, 0},
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
		{nil, &WordList{baseList}, false},
		{&WordList{baseList}, &WordList{baseList}, true},
		{&WordList{baseList}, &WordList{[]string{"foo"}}, false},
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
	cases := []struct {
		filter string
		want   WordList
	}{
		{".oo", WordList{[]string{"bar", "bam"}}},
		{"ba", WordList{[]string{"foo", "zoo"}}},
		{"....", WordList{baseList}},
		{"...", WordList{[]string{}}},
		{"..", WordList{[]string{}}},
		{"", WordList{}},
		{"nomatch", WordList{baseList}},
	}

	for _, c := range cases {
		t.Run(c.filter, func(t *testing.T) {
			w := &WordList{baseList}
			w.Delete(regexp.MustCompile(c.filter))

			if got, want := w.Length(), c.want.Length(); got != want {
				t.Fatalf("got %d words %#v, want %d %#v", got, w, want, c.want)
			}

			if !w.Equals(&c.want) {
				t.Errorf("got %#v, want %#v", w, c.want)
			}
		})
	}
}

func TestKeepOnly(t *testing.T) {
	cases := []struct {
		filter string
		want   WordList
	}{
		{".oo", WordList{[]string{"foo", "zoo"}}},
		{"ba", WordList{[]string{"bar", "bam"}}},
		{"....", WordList{}},
		{"...", WordList{baseList}},
		{"..", WordList{baseList}},
		{"", WordList{baseList}},
		{"nomatch", WordList{}},
	}

	for _, c := range cases {
		t.Run(c.filter, func(t *testing.T) {
			w := &WordList{baseList}
			w.KeepOnly(regexp.MustCompile(c.filter))

			if got, want := w.Length(), c.want.Length(); got != want {
				t.Fatalf("got %d words %#v, want %d %#v", got, w, want, c.want)
			}

			if !w.Equals(&c.want) {
				t.Errorf("got %#v, want %#v", w, c.want)
			}
		})
	}
}
