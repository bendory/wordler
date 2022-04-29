package wordler

import "testing"

func TestDictionary(t *testing.T) {
	d, err := NewDictionary()
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}

	if d.Length() < 10000 {
		t.Errorf("Dictionary looks small: only found %d words.", d.Length())
	}
}
