# wordler
Play and solve wordle!

## WordList
WordList is just a list of words.

## Solver
Solver will solve your wordle for you!

## Puzzler
Puzzler will run a wordle for you.

## Statistics

All stats are based on 1000 6-guess iterations on 5-letter worldles.

* Each guess is simply a random choice:
	* 91.5% success rate
	* Average guesses to win: 4.55 guesses
* Optimize guess based *only* on letter frequency:
	* 79.0% success rate
	* Average guesses to win: 4.91 guesses
	* This is *worse than random* because words with repeated common letters
	  become optimal -- the first guess is always 'arara', which only includes 2
	  letters!


## TODO
* [x] make wordle puzzle, give feedback on guesses
* [x] make solver
    * [ ] optimize solver
* [x] connect puzzler to solver
    * [x] gather statistics on iterations
* [ ] problems: /usr/share/dict does not include...
    * [ ] plurals
	* [ ] fewer
	* [ ] heist
* [ ] load dictionary cross-platform
