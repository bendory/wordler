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

type stats struct {
	iterations, puzzlerFailures, solverFailures, invalidGuesses int
	noWordsRemaining, outOfGuesses, badReactions, winners       int
	winningIteration                                            float32
}

var verbosity int = -1

func main() {
	args := &puzzler.Args{}
	flag.BoolVar(&args.Hard, "hard", true, "use hard rules: 'Any revealed hints must be used in subsequent guesses'")
	flag.IntVar(&args.WordLength, "length", wordler.DEFAULT_WORD_LENGTH, "word length")
	flag.IntVar(&args.Guesses, "guesses", wordler.DEFAULT_GUESSES, "number of guesses allowed")
	flag.StringVar(&args.Solution, "solution", "", "puzzler will use the specified solution")
	iterations := flag.Int("iterations", 10, "number of iterations to run")
	flag.IntVar(&verbosity, "verbosity", verbosity, "-1 (no debug output); 0+ increasing verbosity")
	usage := flag.Usage
	flag.Usage = func() {
		usage()
		fmt.Fprintf(flag.CommandLine.Output(), "\nRemaining positional arguments are taken as guesses to feed to solver.\n")
	}
	flag.Parse()
	clGuesses := flag.Args()

	fmt.Println("I'm a wordler! I try to solve wordle puzzles and report on my success.")
	fmt.Printf("I only allow %d-letter words found in the local dictionary.\n", args.WordLength)
	fmt.Printf("I allow %d guesses for each of %d iterations.\n", args.Guesses, *iterations)
	if args.Solution != "" {
		fmt.Printf("I'll always use '%v' as my solution.\n", args.Solution)
	}
	if len(clGuesses) > 0 {
		fmt.Printf("My first guesses, in order, will be %v.\n", strings.Join(clGuesses, ", "))
		if args.Solution != "" && *iterations != 1 {
			fmt.Println("NOTE: flags set both solution and guesses; setting iterations to 1.")
			*iterations = 1
		}
	}
	fmt.Println("Ready? Here we go!")
	fmt.Println()

	count := stats{}

	option := wordlist.KeepOnlyOption{Exp: regexp.MustCompile(fmt.Sprintf("^.{%d}$", args.WordLength))}
	winningResponse := strings.Repeat(string(wordler.CORRECT), args.WordLength)
	for i := 0; i < *iterations; i++ {
		count.iterations++
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
				fmt.Printf("  ERROR: %d Puzzler words != %d Solver words (continuing anyway)\n", p.Words(), s.Remaining())
			}
			debug(0, "  %d guesses and %d words remain.", p.Guesses(), p.Words())

		GUESS:
			for {
				if len(clGuesses) > 0 {
					guess = clGuesses[0]
					clGuesses = clGuesses[1:]
				} else {
					guess = s.Guess()
				}
				response, err = p.Guess(guess)
				switch err {
				case puzzler.InvalidGuessErr, puzzler.NotInDictionaryErr:
					fmt.Printf("  Invalid guess '%v': %v\n", guess, err)
					count.invalidGuesses++
					s.NotInWordle(guess)
				case puzzler.OutOfGuessesErr:
					count.outOfGuesses++
					break GUESS
				case puzzler.NoWordsRemainingErr:
					fmt.Println("  Uh oh, no words remaining in Puzzler!?")
					count.noWordsRemaining++
					break GUESS
				case nil:
					break GUESS
				}
			}
			if response == winningResponse {
				fmt.Printf("  WINNER! '%v' is the word!\n", guess)
				break
			} else {
				debug(1, "  '%v' --> '%v'", guess, response)
				if err = s.React(guess, response); err != nil {
					count.badReactions++
					fmt.Printf("  ERROR: guess '%v' --> %v\n", guess, err)
				}
			}
		}

		if response == winningResponse {
			count.winners++
			count.winningIteration += float32(args.Guesses - p.Guesses())
		} else if p.Guesses() == 0 {
			fmt.Println("  YOU LOSE!")
		}
		fmt.Printf("  The solution is '%v'.\n", p.GiveUp())
		fmt.Println()
	}

	count.winningIteration /= float32(count.iterations)
	fmt.Printf("Stats gathered: %#v\n", count)
}

// debug prints debug logs
func debug(level int, f string, args ...interface{}) {
	if level <= verbosity {
		fmt.Printf(f, args...)
		fmt.Println()
	}
}
