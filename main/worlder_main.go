package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	"wordler"
	"wordler/puzzler"
	"wordler/solver"
	"wordler/wordlist"
)

func main() {
	args := &puzzler.Args{}
	flag.BoolVar(&args.Hard, "hard", true, "use hard rules: 'Any revealed hints must be used in subsequent guesses'")
	flag.IntVar(&args.WordLength, "length", wordler.DEFAULT_WORD_LENGTH, "word length")
	flag.IntVar(&args.Guesses, "guesses", wordler.DEFAULT_GUESSES, "number of guesses allowed")
	flag.StringVar(&args.Solution, "solution", "", "puzzler will use the specified solution")
	iterations := flag.Int("iterations", 10, "number of iterations to run")
	flag.Parse()

	fmt.Println("I'm a wordler! I try to solve wordle puzzles and report on my success.")
	fmt.Printf("I only allow %d-letter words found in the local dictionary.\n", args.WordLength)
	fmt.Printf("I allow %d guesses for each of %d iterations.\n", args.Guesses, *iterations)
	fmt.Println("Ready? Here we go!")
	fmt.Println()

	option := wordlist.KeepOnlyOption{Exp: regexp.MustCompile(fmt.Sprintf("^.{%d}$", args.WordLength))}
	winningResponse := strings.Repeat(string(wordler.CORRECT), args.WordLength)
	for i := 0; i < *iterations; i++ {
		fmt.Printf("Iteration %d/%d: ", i+1, *iterations)
		p, err := puzzler.New(args)
		if err != nil {
			fmt.Printf("Failed to make a Puzzler: %v\n", err)
			continue
		}
		s, err := solver.New(option)
		if err != nil {
			fmt.Printf("Failed to make a Solver: %v\n", err)
			continue
		}
		fmt.Println()

		var guess, response string
		for p.Guesses() > 0 {
			if p.Words() != s.Remaining() {
				fmt.Printf("  ERROR: %d Puzzler words != %d Solver words (continuing anyway)", p.Words(), s.Remaining())
			}
			fmt.Printf("  %d guesses and %d words remain.\n", p.Guesses(), p.Words())

		GUESS:
			for {
				guess = s.Guess()
				response, err = p.Guess(guess)
				switch err {
				case puzzler.InvalidGuessErr, puzzler.NotInDictionaryErr:
					fmt.Printf("  Invalid guess '%v': %v\n", guess, err)
					s.NotInWordle(guess)
				case puzzler.OutOfGuessesErr:
					break GUESS
				case puzzler.NoWordsRemainingErr:
					fmt.Println("  Uh oh, no words remaining in Puzzler!?")
					break GUESS
				case nil:
					break GUESS
				}
			}
			if response == winningResponse {
				fmt.Printf("  WINNER! '%v' is the word!\n", guess)
				break
			} else {
				fmt.Printf("  '%v' --> '%v'\n", guess, response)
				if err = s.React(guess, response); err != nil {
					fmt.Printf("  ERROR: guess '%v' --> %v\n", guess, err)
				}
			}
		}

		if p.Guesses() == 0 && response != winningResponse {
			fmt.Println("YOU LOSE!")
		}
		fmt.Printf("  The solution is '%v'.\n", p.GiveUp())
		fmt.Println()
	}
}
