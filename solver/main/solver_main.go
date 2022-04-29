package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"wordler"
	"wordler/solver"
	"wordler/wordlist"
)

func main() {
	length := flag.Int("length", wordler.DEFAULT_WORD_LENGTH, "word length")
	usage := flag.Usage
	flag.Usage = func() {
		usage()
		fmt.Fprintf(flag.CommandLine.Output(), "\nRemaining positional arguments are taken as guesses to feed to solver.\n")
	}
	flag.Parse()

	fmt.Println("I'm a wordle solver! I'll make guesses, you tell me wordle's response.")
	fmt.Printf("I only allow %d-letter words found in the local dictionary.\n", *length)
	fmt.Printf("Use '%c' for \"right letter in the right place\"\n", wordler.CORRECT)
	fmt.Printf("Use '%c' for \"right letter in the wrong place\"\n", wordler.ELSEWHERE)
	fmt.Printf("Use '%c' for \"letter not in the word\"\n", wordler.NIL)
	fmt.Println("Respond with the letter 'n' by itself to tell me that my guess isn't in wordle's dictionary.")
	fmt.Println("Ready? Here we go!")
	fmt.Println()

	g, err := solver.New(wordlist.KeepOnlyOption{Exp: regexp.MustCompile(fmt.Sprintf("^.{%d}$", *length))})

	if err != nil {
		fmt.Printf("Failed to make a Solver: %v\n", err)
		os.Exit(2)
	}

	clGuesses := flag.Args()
	for {
		switch g.Remaining() {
		case 0:
			fmt.Println("ERROR: solver is empty.")
			os.Exit(1)

		case 1:
			fmt.Println("The word is " + g.Guess())
			os.Exit(0)

		default:
			fmt.Printf("I've got %d possible words left.\n", g.Remaining())

			var guess string
			if len(clGuesses) > 0 {
				guess = clGuesses[0]
				clGuesses = clGuesses[1:]
			} else {
				guess = g.Guess()
			}
			fmt.Println("Guess: " + guess)

			for done := false; !done; {
				var response string
				fmt.Print("Response? ")
				fmt.Scan(&response)

				if response == "n" {
					g.NotInWordle(guess)
					done = true
				} else if err := g.React(guess, response); err != nil {
					fmt.Println("ERROR: ", err)
					fmt.Printf("Guess was \"%v\"\n", guess)
				} else {
					done = true
				}
			}
			fmt.Println()
		}
	}
}
