package wordlist

import (
	"regexp"
	"testing"
)

func TestLength(t *testing.T) {
	baseList := []string{"foo", "bar", "bam", "zoo"}
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
	baseList := []string{"foo", "bar", "bam", "zoo"}
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
	baseList := []string{"foo", "bar", "bam", "zoo"}
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
	baseList := []string{"foo", "bar", "bam", "zoo"}
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
	baseList := []string{"foo", "bar", "bam", "zoo"}
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

// TestPlatformDictionary is non-hermetic; it loads the actual dictionary for
// this platform.
func TestPlatformDictionary(t *testing.T) {
	d, err := NewDictionary()
	if err != nil {
		t.Errorf("Failed to load dictionary: %v", err)
	}
	if d.Length() < 10000 {
		t.Errorf("Dictionary looks small: only found %d words.", d.Length())
	}

	d, err = NewDictionary(KeepOnlyOption{regexp.MustCompile("^i$")})
	if err != nil {
		t.Errorf("Failed to load dictionary: %v", err)
	}
	if d.Length() != 1 {
		t.Errorf("want 1 word, got %d", d.Length())
	}
	want := New([]string{"i"})
	if !want.Equals(d) {
		t.Errorf("want %#v != got %#v", want, d)
	}
}

func TestContains(t *testing.T) {
	baseList := []string{"foo", "bar", "bam", "zoo"}
	w := New(baseList)
	if got, want := w.Length(), len(baseList); got != want {
		t.Errorf("got %d; want %d", got, want)
	}

	for _, word := range baseList {
		if !w.Contains(word) {
			t.Errorf("Dictionary does not contain %v.", word)
		}
	}
	if w.Contains("not a word") {
		t.Errorf("Dictionary contains 'not a word'.")
	}
	if w.Contains("") {
		t.Errorf("Dictionary contains empty string.")
	}

	// Don't error on nil.
	w = nil
	if w.Contains("foo") {
		t.Errorf("nil list contains foo")
	}
	if w.Length() != 0 {
		t.Errorf("nil list has %d words", w.Length())
	}
}

func TestOptions(t *testing.T) {
	baseList := []string{"a", "ab", "abc"}
	cases := []struct {
		desc     string
		baseList []string
		options  []Option
		want     *WordList
	}{
		{"K1", baseList, []Option{KeepOnlyOption{regexp.MustCompile("^..$")}}, New([]string{"ab"})},
		{"D1", baseList, []Option{DeleteOption{regexp.MustCompile("^.$")}}, New([]string{"ab", "abc"})},
		{"D2", baseList, []Option{DeleteOption{regexp.MustCompile("b")}}, New([]string{"a"})},
		{"K2", baseList, []Option{KeepOnlyOption{regexp.MustCompile("c$")}}, New([]string{"abc"})},
		{
			"KD",
			baseList,
			[]Option{KeepOnlyOption{regexp.MustCompile("..")}, DeleteOption{regexp.MustCompile("c$")}},
			New([]string{"ab"}),
		}, {
			"DD",
			baseList,
			[]Option{DeleteOption{regexp.MustCompile("^..$")}, DeleteOption{regexp.MustCompile("^a$")}},
			New([]string{"abc"}),
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			got := New(c.baseList, c.options...)
			if !got.Equals(c.want) {
				t.Errorf("want %#v; got %#v", c.want, got)
			}
		})
	}
}

func TestOptimalGuess(t *testing.T) {
	cases := []struct {
		list []string
		want string
	}{
		{[]string{"aaa", "bcd"}, "bcd"},
		{[]string{"bcd", "aaa"}, "bcd"},
		{[]string{"aaa", "bcd", "def", "hij", "cic", "ccc"}, "bcd"},
	}

	for _, c := range cases {
		t.Run(c.want, func(t *testing.T) {
			w := New(c.list)
			if want, got := c.want, w.OptimalGuess(); want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		})
	}
}
