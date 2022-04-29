package main

import (
	"fmt"
	"os"
	"regexp"

	"wordler/guesser"
)

func main() {
	fmt.Println("I'm a wordle guesser! I'll make guesses, you tell me wordle's response.")
	fmt.Printf("Use '%c' for \"right letter in the right place\"\n", guesser.CORRECT)
	fmt.Printf("Use '%c' for \"right letter in the wrong place\"\n", guesser.ELSEWHERE)
	fmt.Printf("Use '%c' for \"letter not in the word\"\n", guesser.NIL)
	fmt.Println("Respond with the letter 'n' by itself to tell me that my guess isn't in wordle's dictionary.")
	fmt.Println("Ready? Here we go!")
	fmt.Println()

	g, err := guesser.New(guesser.KeepOnlyOption{regexp.MustCompile("^.....$")})
	if err != nil {
		fmt.Printf("Failed to make a Guesser: %v\n", err)
		os.Exit(2)
	}
	
	for {
		switch g.Remaining() {
		case 0:
			fmt.Println("ERROR: guesser is empty.")
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
