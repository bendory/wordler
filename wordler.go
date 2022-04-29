package wordler

const (
	// WORD_LENGTH is the puzzle size in letters.
	WORD_LENGTH = 5

	// CORRECT indicates that the letter is in the puzzle in the given location.
	CORRECT = '+'

	// ELSEWHERE indicates that the letter is somewhere else in the puzzle.
	ELSEWHERE = '*'

	// NIL indicates that the letter is not found in the puzzle.
	NIL = '_'
)
