package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"wordler"
	"wordler/solver"
	"wordler/wordlist"
)

func main() {
	local := flag.Bool("local_dictionary", false, "use local dictionary in place of Wordle dictionary")
	length := flag.Int("length", wordler.DEFAULT_WORD_LENGTH, "word length")
	guesses := flag.Uint("guesses", wordler.DEFAULT_GUESSES, "number of guesses allowed")
	usage := flag.Usage
	flag.Usage = func() {
		usage()
		fmt.Fprintf(flag.CommandLine.Output(), "\nRemaining positional arguments are taken as guesses to feed to solver.\n")
	}
	flag.Parse()

	fmt.Printf("I'm a wordle solver! I'll make up to %d guesses, you tell me wordle's response.\n", *guesses)
	if *local {
		fmt.Printf("I only allow %d-letter words found in the local dictionary.\n", *length)
	} else if *length != wordler.DEFAULT_WORD_LENGTH && *length != 0 {
		fmt.Printf("I'm ignoring the --length=%d flag, using Worlde dictionary\n", *length)
	}
	fmt.Printf("Use '%c' for \"right letter in the right place\"\n", wordler.CORRECT)
	fmt.Printf("Use '%c' for \"right letter in the wrong place\"\n", wordler.ELSEWHERE)
	fmt.Printf("Use '%c' for \"letter not in the word\"\n", wordler.NIL)
	fmt.Println("Respond with the letter 'n' by itself to tell me that my guess isn't in wordle's dictionary.")
	fmt.Println("Respond with the letter 'y' by itself to tell me that I've solved the wordle.")
	fmt.Println("Ready? Here we go!")
	fmt.Println()

	var s *solver.Solver
	var err error
	if *local {
		s, err = solver.New(wordlist.KeepOnlyOption{Exp: regexp.MustCompile(fmt.Sprintf("^.{%d}$", *length))})
	} else {
		s = solver.From(wordler.Dictionary)
	}
	if err != nil {
		fmt.Printf("Failed to make a Solver: %v\n", err)
		os.Exit(2)
	}

	clGuesses := flag.Args()
	for ; *guesses > 0; *guesses-- {
		switch s.Remaining() {
		case 0:
			fmt.Println("ERROR: solver is empty.")
			os.Exit(1)

		case 1:
			fmt.Println("The word is " + s.Guess())
			os.Exit(0)

		default:
			fmt.Printf("I've got %d possible words and %d guesses left.\n", s.Remaining(), *guesses)

			var guess string
			if len(clGuesses) > 0 {
				guess = clGuesses[0]
				clGuesses = clGuesses[1:]
			} else {
				guess = s.Guess()
			}
			fmt.Println("Guess: " + guess)

			for done := false; !done; {
				var response string
				fmt.Print("Response? ")
				fmt.Scan(&response)

				switch response {
				case "n":
					s.NotInWordle(guess)
					done = true
					continue

				case "y": // syntactic sugar for all-correct response
					response = strings.Repeat(string(wordler.CORRECT), *length)
				}

				if err := s.React(guess, response); err != nil {
					fmt.Println("ERROR: ", err)
					fmt.Printf("Guess was \"%v\"\n", guess)
				} else {
					done = true
				}
			}
			fmt.Println()
		}
	}
	fmt.Println("Out of guesses :-(")
	os.Exit(0)
}
