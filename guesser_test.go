package wordler

import "testing"

var testList = []string{"foo", "bar", "bam", "boom"}

func TestNewGuesser(t *testing.T) {
	g, err := NewGuesser()
	if err != nil {
		t.Fatal(err)
	}

	if g.w.Length() == 0 {
		t.Errorf("Only found %d words in initial list.", g.w.Length())
	}
}

func TestGuess(t *testing.T) {
	guesser := &Guesser{w: NewWordList(testList)}
	guess := guesser.Guess()

	if !guesser.w.Contains(guess) {
		t.Errorf("Guess %v not found in wordlist %#v", guess, guesser.w)
	}

	singleton := "foo"
	guesser.w = NewWordList([]string{singleton})
	guess = guesser.Guess()
	if guess != singleton {
		t.Errorf("Want guess %v, got %v", singleton, guess)
	}

	guesser.w = nil
	guess = guesser.Guess()
	if guess != "" {
		t.Errorf("Want empty string, got %v", guess)
	}
}
