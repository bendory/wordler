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
