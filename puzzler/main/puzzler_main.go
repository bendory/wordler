package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"wordler"
	"wordler/puzzler"
)

func main() {
	args := &puzzler.Args{}
	flag.BoolVar(&args.Hard, "hard", true, "use hard rules: 'Any revealed hints must be used in subsequent guesses'")
	flag.IntVar(&args.WordLength, "length", 5, "word length")
	flag.IntVar(&args.Guesses, "guesses", 6, "number of guesses allowed")
	flag.StringVar(&args.Solution, "solution", "", "puzzler will use the specified solution")
	flag.Parse()

	fmt.Println("I'm a wordle puzzle! You make guesses, I'll score them.")
	fmt.Printf("I only allow %d-letter words found in the local dictionary.\n", wordler.DEFAULT_WORD_LENGTH)
	fmt.Printf("I'll use '%c' for \"right letter in the right place\"\n", wordler.CORRECT)
	fmt.Printf("I'll use '%c' for \"right letter in the wrong place\"\n", wordler.ELSEWHERE)
	fmt.Printf("I'll use '%c' for \"letter not in the word\"\n", wordler.NIL)
	fmt.Println("I'll respond with the letter 'n' by itself if your guess isn't in wordle's dictionary.")
	fmt.Println("Ready? Here we go!")
	fmt.Println()

	p, err := puzzler.New(args)

	if err != nil {
		fmt.Printf("Failed to make a Puzzler: %v\n", err)
		os.Exit(2)
	}

	winningResponse := strings.Repeat(string(wordler.CORRECT), args.WordLength)
	for p.Guesses() > 0 {
		fmt.Printf("%d guesses and %d words remain.\n", p.Guesses(), p.Words())
		var guess, response string

	GUESS:
		for {
			fmt.Print("Your guess? ")
			fmt.Scan(&guess)

			var err error
			response, err = p.Guess(guess)
			switch err {
			case puzzler.InvalidGuessErr, puzzler.NotInDictionaryErr:
				fmt.Println("Invalid guess: ", err)
			case puzzler.OutOfGuessesErr:
				break GUESS
			case puzzler.NoWordsRemainingErr:
				fmt.Println("Uh oh, no words remaining!?")
				break GUESS
			case nil:
				break GUESS
			}
		}
		if response == winningResponse {
			fmt.Println("YOU WIN!")
			break
		} else {
			fmt.Println("Response:  ", response)
			fmt.Println()
		}
	}

	if p.Guesses() == 0 {
		fmt.Println("YOU LOSE!")
	}
	fmt.Printf("The solution is '%v'.\n", p.GiveUp())
}
