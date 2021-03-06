package puzzler

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"wordler"
	"wordler/wordlist"
)

type Dictionary int

const (
	// WordleDictionary signals use of Wordle's dictionary.
	WordleDictionary = 0
	// LocalDictionary signals use of the platform's local dictionary.
	LocalDictionary = 1
)

// Wordle is a Wordle puzzle.
type Wordle struct {
	dict             *wordlist.WordList // full dictionary
	remaining        *wordlist.WordList // words remaining
	word             string             // the answer
	remainingGuesses int                // how many guesses are left
	hard             bool               // hard or easy rules?
}

// Args are used to construct a new Wordle puzzle.
type Args struct {
	Dictionary          Dictionary // dictionary to use; either Wordle or Local
	Hard                bool       // hard rules
	WordLength, Guesses int
	Solution            string // create a puzzler with this solution; otherwise a random word is chosen
	Options             []wordlist.Option
}

var (
	InvalidGuessErr     = errors.New("invalid guess")
	NotInDictionaryErr  = errors.New("not in dictionary")
	NoWordsRemainingErr = errors.New("no words remaining")
	OutOfGuessesErr     = errors.New("no remaining guesses")
	verbose             = false
)

// New creates a new Wordle puzzle, limiting allowed words based on given
// options.
func New(a *Args) (*Wordle, error) {
	if a == nil {
		a = &Args{
			Hard:       true,
			WordLength: wordler.DEFAULT_WORD_LENGTH,
			Guesses:    wordler.DEFAULT_GUESSES,
		}
	}

	w := &Wordle{remainingGuesses: a.Guesses, hard: a.Hard}
	var err error
	switch a.Dictionary {
	case WordleDictionary:
		switch a.WordLength {
		case wordler.DEFAULT_WORD_LENGTH, 0:
			// Allowed
		default:
			return nil, fmt.Errorf("invalid args: cannot specify word length != %d when using Wordle dictionary", wordler.DEFAULT_WORD_LENGTH)
		}
		w.dict = wordlist.New(wordler.Dictionary, a.Options...)
		w.remaining = wordlist.New(wordler.Dictionary, a.Options...)

	case LocalDictionary:
		a.Options = append(a.Options, wordlist.KeepOnlyOption{Exp: regexp.MustCompile(fmt.Sprintf("^[a-z]{%d}$", a.WordLength))})

		if w.dict, err = wordlist.NewDictionary(a.Options...); err != nil {
			return nil, err
		}
		if w.remaining, err = wordlist.NewDictionary(a.Options...); err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("invalid dictionary option")
	}

	if w.Words() == 0 {
		return nil, NoWordsRemainingErr
	}
	if a.Solution != "" {
		if err := w.validate(a.Solution); err == nil {
			w.word = a.Solution
		} else {
			return nil, err
		}
	} else {
		w.word = w.dict.Random()
	}
	return w, nil
}

// Guess the given word.
// The returned string is populated with wordler.CORRECT, wordler.NIL,
// wordler.ELSEWHERE corresponding to the guess.
func (w *Wordle) Guess(g string) (string, error) {
	if w == nil || w.remainingGuesses == 0 {
		return "", OutOfGuessesErr
	}
	if err := w.validate(g); err != nil {
		return "", err
	}
	w.remainingGuesses--

	word := w.word

	// To score the guess, we:
	// - first score any letters as CORRECT
	// - next score letters as ELSEWHERE
	// - remaining letters are scored NIL
	// As we score, we transform guess g into the returned score.
	// As we score, we delete letters from word to prevent double-scoring.
	lettersFound := make(map[byte]bool)

	// First score all the letters in the CORRECT place.
	for i, b := range g {
		c := byte(b)
		if word[i] == c {
			g = g[:i] + string(wordler.CORRECT) + g[i+1:]
			word = word[:i] + string(wordler.CORRECT) + word[i+1:] // prevent additional matches on this letter
			w.remaining.KeepOnly(regexp.MustCompile(fmt.Sprintf("^%s%s", strings.Repeat(".", i), string(c))))
			lettersFound[c] = true
		}
	}
	debug("after scoring %c, guess is %v, word is %v", wordler.CORRECT, g, word)

	// Now score all the letters that appear ELSEWHERE in word.
	for i, b := range g { // i is the letter's position in g
		c := byte(b)
		if c == wordler.CORRECT {
			continue
		}

		for j, b := range word { // j is the position we are comparing in word
			l := byte(b)
			if j == i || l == wordler.CORRECT || l == wordler.ELSEWHERE {
				continue
			}
			if l == c { // if l in word matches c in guess
				g = g[:i] + string(wordler.ELSEWHERE) + g[i+1:]
				word = word[:j] + string(wordler.ELSEWHERE) + word[j+1:] // prevent additional matches on this letter
				debug("keeping only words containing '%c'", c)
				w.remaining.KeepOnly(regexp.MustCompile(string(c)))
				debug("%d words left.", w.remaining.Length())
				pattern := fmt.Sprintf("^%s%s", strings.Repeat(".", i), string(c))
				debug("deleting all words with '%c' as char %d by using regexp '%v'", c, i+1, pattern)
				w.remaining.Delete(regexp.MustCompile(pattern))
				debug("%d words left.", w.remaining.Length())
				lettersFound[c] = true
				break
			}
		}

		if g[i] != wordler.ELSEWHERE {
			g = g[:i] + string(wordler.NIL) + g[i+1:]
			if !lettersFound[c] {
				debug("deleting all words containing '%c'", c)
				w.remaining.Delete(regexp.MustCompile(string(c)))
			} else {
				pattern := fmt.Sprintf("^%s%s", strings.Repeat(".", i), string(c))
				debug("deleting all words with '%c' as char %d by using regexp '%v'", c, i+1, pattern)
				w.remaining.Delete(regexp.MustCompile(pattern))
			}
			debug("%d words left.", w.remaining.Length())
		}
	}

	// g is now a combination of CORRECT, ELSEWHERE, and NIL
	return g, nil
}

// validate guess based on `hard` setting.
func (w *Wordle) validate(g string) error {
	if !w.dict.Contains(g) {
		return fmt.Errorf("'%s' %w", g, NotInDictionaryErr)
	}
	if w.Words() == 0 {
		return NoWordsRemainingErr
	}
	if w.hard && !w.remaining.Contains(g) {
		return InvalidGuessErr
	}
	return nil
}

// Guesses returns the number of guesses left.
func (w *Wordle) Guesses() int {
	if w == nil {
		return 0
	}
	return w.remainingGuesses
}

// Words returns the number of words remaining. The wordle puzzle tracks words
// remaining assuming you apply all guess responses correctly and only guess
// using `hard` rules.
func (w *Wordle) Words() int {
	if w == nil {
		return 0
	}
	return w.remaining.Length()
}

// GiveUp: no more guesses are allowed and the solution is revealed.
func (w *Wordle) GiveUp() string {
	if w == nil {
		return ""
	}
	w.remainingGuesses = 0
	w.remaining = wordlist.New([]string{w.word})
	return w.word
}

// debug prints debug logs
func debug(f string, args ...interface{}) {
	if verbose {
		fmt.Printf(f, args...)
		fmt.Println()
	}
}
