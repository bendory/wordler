package wordler

import (
	"regexp"
	"testing"
)

func TestFilter(t *testing.T) {
	w := WordList{"foo", "bar", "bam", "zoo"}

	cases := []struct {
		filter string
		want   WordList
	}{
		{".oo", []string{"bam", "bar"}},
		{"ba", []string{"foo", "zoo"}},
		{"....", w},
		{"...", []string{}},
	}

	for _, c := range cases {
		t.Run(c.filter, func(t *testing.T) {
			r := regexp.MustCompile(c.filter)
			got := w.Filter(r)

			if gotLen, wantLen := len(got), len(c.want); gotLen != wantLen {
				t.Fatalf("%d != %d: got %#v, want %#v", gotLen, wantLen, got, c.want)
			}
			for i, gotItem := range got {
				if wantItem := c.want[i]; gotItem != wantItem {
					t.Errorf("[%d]: got %v, want %v", i, gotItem, wantItem)
				}
			}
		})
	}
}
