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
	verbose            = false
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
	r := &Response{GuessesRemaining: w.remainingGuesses}

	word := w.word

	// To score the guess, we:
	// - first score any letters as CORRECT
	// - next score letters as ELSEWHERE
	// - remaining letters are scored NIL
	// As we score, we transform guess g into the returned score.
	// As we score, we delete letters from word to prevent double-scoring.

	// First score all the letters in the CORRECT place.
	for i, b := range g {
		c := byte(b)
		if word[i] == c {
			g = g[:i] + string(wordler.CORRECT) + g[i+1:]
			word = word[:] + string(wordler.CORRECT) + word[i+1:] // prevent additional matches on this letter
			w.remaining.KeepOnly(regexp.MustCompile(fmt.Sprintf("^%s%s", strings.Repeat(".", i), string(c))))
		}
	}

	// Now score all the letters that appear ELSEWHERE in word.
	for i, b := range g {
		c := byte(b)
		if c == wordler.CORRECT {
			continue
		}

		for j, b := range word {
			l := byte(b)
			if j == i || l == wordler.CORRECT || l == wordler.ELSEWHERE {
				continue
			}
			if l == c {
				g = g[:i] + string(wordler.ELSEWHERE) + g[i+1:]
				word = word[:j] + string(wordler.ELSEWHERE) + word[j+1:] // prevent additional matches on this letter
				debug("keeping only words containing '%c'", c)
				w.remaining.KeepOnly(regexp.MustCompile(string(c)))
				debug("%d words left.", w.remaining.Length())
				debug("deleting all words with '%c' as char %d", c, j)
				w.remaining.Delete(regexp.MustCompile(fmt.Sprintf("^%s[^%s]", strings.Repeat(".", j), string(c))))
				debug("%d words left.", w.remaining.Length())
				break
			}
		}

		if g[i] != wordler.ELSEWHERE {
			g = g[:i] + string(wordler.NIL) + g[i+1:]
			debug("deleting all words containing '%c'", c)
			w.remaining.Delete(regexp.MustCompile(string(c)))
			debug("%d words left.", w.remaining.Length())
		}
	}

	// g is now a combination of CORRECT, ELSEWHERE, and NIL
	r.Detail = g
	r.WordsRemaining = w.remaining.Length()
	return r, nil
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

// debug prints debug logs
func debug(f string, args ...interface{}) {
	if verbose {
		fmt.Printf(f, args...)
		fmt.Println()
	}
}
