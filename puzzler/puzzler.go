package wordler

import (
	"errors"

	"wordler/wordlist"
)

const defaultGuesses = 6

// Wordle is a Wordle puzzle.
type Wordle struct {
	dict             *wordlist.WordList // full dictionary
	remaining        *wordlist.WordList // words remaining
	word             string             // the answer
	remainingGuesses int                // how many guesses are left
	strict           bool               // strict or lenient?
}

var (
	NotInDictionaryErr = errors.New("not in dictionary")
	InvalidGuessErr    = errors.New("invalid guess")
)

// New creates a new Wordle puzzle, limiting allowed words based on given
// options.
func New(strict bool, options ...wordlist.Option) (*Wordle, error) {
	var err error
	w := &Wordle{remainingGuesses: defaultGuesses, strict: strict}
	if w.dict, err = wordlist.NewDictionary(options...); err != nil {
		return nil, err
	}
	if w.remaining, err = wordlist.NewDictionary(options...); err != nil {
		return nil, err
	}
	w.word = w.dict.Random()
	return w, nil
}

// Guess the given word.
func (w *Wordle) Guess(g string) (wordsRemaining, guessesRemaining int, err error) {
	if err = w.validate(g); err != nil {
		return w.remaining.Length(), w.remainingGuesses, err
	}
	// TODO: evaluate guess
	// TODO: reduce w.remaining
	w.remainingGuesses--
	return w.remaining.Length(), w.remainingGuesses, err
}

// validate guess based on strictness setting.
func (w *Wordle) validate(g string) error {
	if !w.dict.Contains(g) {
		return NotInDictionaryErr
	}
	if w.strict && !w.remaining.Contains(g) {
		return InvalidGuessErr
	}
	return nil
}
