package main

import (
	"fmt"
	"os"
	"regexp"

	"wordler/solver"
)

func main() {
	fmt.Println("I'm a wordle solver! I'll make guesses, you tell me wordle's response.")
	fmt.Printf("Use '%c' for \"right letter in the right place\"\n", solver.CORRECT)
	fmt.Printf("Use '%c' for \"right letter in the wrong place\"\n", solver.ELSEWHERE)
	fmt.Printf("Use '%c' for \"letter not in the word\"\n", solver.NIL)
	fmt.Println("Respond with the letter 'n' by itself to tell me that my guess isn't in wordle's dictionary.")
	fmt.Println("Ready? Here we go!")
	fmt.Println()

	g, err := solver.New(solver.KeepOnlyOption{Exp: regexp.MustCompile("^.....$")})
	if err != nil {
		fmt.Printf("Failed to make a Solver: %v\n", err)
		os.Exit(2)
	}

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
			var response string
			guess := g.Guess()
			fmt.Println("Guess: " + guess)
			fmt.Print("Response? ")
			fmt.Scan(&response)
			if response != "n" {
				g.React(guess, response)
			}
			fmt.Println()
		}
	}
}