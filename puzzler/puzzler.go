package wordler

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"wordler"
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

type Response struct {
	Detail                           string
	WordsRemaining, GuessesRemaining int
}

// Guess the given word.
func (w *Wordle) Guess(g string) (*Response, error) {
	if err := w.validate(g); err != nil {
		return nil, err
	}
	w.remainingGuesses--
	r := &Response{
		GuessesRemaining: w.remainingGuesses,
		Detail:           w.evaluateGuess(g),
	}
	r.WordsRemaining = w.remaining.Length()

	return r, nil
}

func (w *Wordle) evaluateGuess(g string) string {
	var response []rune
	word := w.word
	for i, b := range g {
		c := byte(b)
		r := wordler.NIL
		if word[i] == c {
			r = wordler.CORRECT
			w.remaining.KeepOnly(regexp.MustCompile(fmt.Sprintf("^%s%s", strings.Repeat(".", i), string(c))))
		} else {
			for j := i + 1; j < len(word); j++ {
				if word[j] == c && g[j] != c {
					r = wordler.ELSEWHERE
					word = word[:j] + " " + word[j+1:] // prevent additional matches on this letter
					w.remaining.KeepOnly(regexp.MustCompile(string(c)))
					w.remaining.Delete(regexp.MustCompile(fmt.Sprintf("^%s[^%s]", strings.Repeat(".", i), string(c))))
					break
				}
			}

			if r == wordler.NIL {
				w.remaining.Delete(regexp.MustCompile(string(c)))
			}
		}
		response = append(response, r)
	}
	return string(response)
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
