package main

import (
	"flag"
	"fmt"
	"os"
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
	local := flag.Bool("local_dictionary", false, "use local dictionary in place of Wordle dictionary")
	flag.IntVar(&args.WordLength, "length", wordler.DEFAULT_WORD_LENGTH, "word length")
	flag.IntVar(&args.Guesses, "guesses", wordler.DEFAULT_GUESSES, "number of guesses allowed")
	flag.StringVar(&args.Solution, "solution", "", "puzzler will use the specified solution")
	iterations := flag.Int("iterations", 10, "number of iterations to run")
	flag.IntVar(&verbosity, "verbosity", verbosity, "-2 (silent); -1 (no debug output); 0+ increasing verbosity")
	usage := flag.Usage
	flag.Usage = func() {
		usage()
		fmt.Fprintf(flag.CommandLine.Output(), "\nRemaining positional arguments are taken as guesses to feed to solver.\n")
	}
	flag.Parse()
	clGuesses := flag.Args()

	fmt.Println("I'm a wordler! I try to solve wordle puzzles and report on my success.")
	if *local {
		fmt.Printf("I only allow %d-letter words found in the local dictionary.\n", args.WordLength)
		args.Dictionary = puzzler.LocalDictionary
	}
	fmt.Printf("I allow %d guesses for each of %d iterations.\n", args.Guesses, *iterations)
	if args.Solution != "" {
		fmt.Printf("I'll always use '%v' as my solution.\n", args.Solution)
	}
	if len(clGuesses) > 0 {
		if *iterations != 1 && clGuesses[len(clGuesses)-1] == args.Solution {
			fmt.Println("NOTE: last guess is solution; setting iterations to 1.")
			*iterations = 1
		} else {
			fmt.Printf("My first guesses, in order, will be %v.\n", strings.Join(clGuesses, ", "))
		}
	}
	fmt.Println("Ready? Here we go!")
	fmt.Println()

	count := stats{}

	option := wordlist.KeepOnlyOption{Exp: regexp.MustCompile(fmt.Sprintf("^.{%d}$", args.WordLength))}
	winningResponse := strings.Repeat(string(wordler.CORRECT), args.WordLength)
	for i := 0; i < *iterations; i++ {
		count.iterations++
		debug(-1, "Iteration %d/%d: ", i+1, *iterations)
		p, err := puzzler.New(args)
		if err != nil {
			count.puzzlerFailures++
			fmt.Printf("Failed to make a Puzzler: %v\n", err)
			os.Exit(1) // This should never happen.
		}

		var s *solver.Solver
		if *local {
			if s, err = solver.New(option); err != nil {
				count.solverFailures++
				fmt.Printf("Failed to make a Solver: %v\n", err)
				os.Exit(1) // This should never happen.
			}
		} else {
			s = solver.From(wordler.Dictionary)
		}

		var guess, response string
		var guesses []string
	OUTER: // Loop until we win, get an error, or run out of guesses.
		for p.Guesses() > 0 {
			if p.Words() != s.Remaining() {
				fmt.Printf("  ERROR: %d Puzzler words != %d Solver words (continuing anyway)\n", p.Words(), s.Remaining())
			}
			debug(0, "  %d guesses and %d words remain.", p.Guesses(), p.Words())

		GUESS: // Loop until we get a valid guess.
			for {
				// Exhaust guesses specified on the command line, then use our
				// guessing algorithm.
				if len(clGuesses) > 0 {
					guess = clGuesses[0]
					clGuesses = clGuesses[1:]
				} else {
					guess = s.Guess()
				}
				guesses = append(guesses, guess)
				response, err = p.Guess(guess)
				switch err {

				// This should never happen given that puzzler and solver use
				// the same dictionary.
				case puzzler.InvalidGuessErr, puzzler.NotInDictionaryErr:
					fmt.Printf("  Invalid guess '%v': %v\n", guess, err)
					count.invalidGuesses++
					s.NotInWordle(guess)

				// This should never happen; we should break out of OUTER before
				// getting this error.
				case puzzler.OutOfGuessesErr:
					count.outOfGuesses++
					break OUTER

				// This should never happen; we should either run out of guesses
				// or win first.
				case puzzler.NoWordsRemainingErr:
					fmt.Println("  Uh oh, no words remaining in Puzzler!?")
					count.noWordsRemaining++
					break OUTER

				// Expected behavior -- valid guess.
				case nil:
					break GUESS
				}
			}

			if response == winningResponse {
				break
			}
			debug(1, "  '%v' --> '%v'", guess, response)
			if err = s.React(guess, response); err != nil {
				count.badReactions++
				fmt.Printf("  ERROR: guess '%v' --> %v\n", guess, err)
			}
		}

		if response == winningResponse {
			debug(-1, "  WINNER! '%v' is the word! Guesses: %v", guess, strings.Join(guesses, ", "))
			count.winners++
			count.winningIteration += float32(args.Guesses - p.Guesses())
		} else if p.Guesses() == 0 {
			debug(-1, "  YOU LOSE!")
			debug(-1, "  Guesses were: %v; %d words left.", strings.Join(guesses, ", "), s.Remaining())
		}
		debug(-1, "  The solution is '%v'.", p.GiveUp())
		debug(-1, "")
	}

	count.winningIteration /= float32(count.winners)
	debug(-1, "Stats gathered: %#v", count)
	fmt.Printf("I won %.2f%% of games played with an average of %.2f guesses.\n",
		float32(count.winners*100)/float32(count.iterations), count.winningIteration)
}

// debug prints debug logs
func debug(level int, f string, args ...interface{}) {
	if level <= verbosity {
		fmt.Printf(f, args...)
		fmt.Println()
	}
}
