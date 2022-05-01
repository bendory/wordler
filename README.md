# wordler
Play and solve wordle!

## WordList
WordList is just a list of words.

By default, it loads the dictionary from `/usr/share/dict/words`, but var
`wordler.Loader` can be changed to load any wordlist.

`wordlist.OptimalGuess()` contains the exciting heuristic to choose the best
next guess.

## Solver
Solver will solve your wordle for you!

The interactive solver can be used to solve any wordle; args allow changing word
length. It always plays using the "hard" rules.

## Puzzler
Puzzler will run a wordle for you; solve it yourself.

Args allow changing word length, number of guesses, and more.

## Main
`wordler/main` connects a Solver to a Puzzler and runs simulated wordle
interactions; it's helpful for gathering statistics on solution success rate.

## Simulator
Simulator is for testing.  It confirms that Solver and Puzzler score guesses and
use them for solving with a reciprocal approach.

## Statistics
All stats are based on 1000 6-guess iterations on 5-letter wordles.

* Using `/usr/share/dict`:
	* Random guesses:
		* 91.5% success rate
		* Average guesses to win: 4.55 guesses
	* Optimize guess based *only* on heaviest weighted-averge letter frequency:
		* 79.0% success rate
		* Average guesses to win: 4.91 guesses
		* This is *worse than random* because words with repeated common letters
		  become optimal -- the first guess is always 'arara', which only includes 2
		  letters! (The second guess was often 'neese' which is also awful.)
	* Optimize based on most-new-letters-in-guess followed by heaviest
	  weighted-average letter frequency:
		* 96.0% success rate
		* Average guesses to win: 4.21 guesses
* Using Wordler dictionary:
	* Random guesses:
		* 84.3% success rate
		* Average guesses to win: 4.58
	* Optimize based on most-new-letters-in-guess followed by heaviest
	  weighted-average letter frequency:
		* 87.0% success rate
		* Average guesses to win: 4.36 guesses
	* Weight letters based on how many words they appear in instead of total
	  number of times they appear (which is above "letter frequency"):
		* 87.9% success rate
		* Average guesses to win: 4.36 guesses

## TODO
* [x] make wordle puzzle, give feedback on guesses
* [x] make solver
    * [x] optimize solver
	* [ ] Optimize further! `egrep ^.ater$ /usr/share/dict/words` reveals a
	  pathological case where `bater` can take 8 guesses. Optimized solver
	  should seek to avoid these pathological cases by avoiding guesses that
	  use a pattern of letters repeated in other words. Another test case:
	  `baker` is solved ~34% of the time with the current algorithm due to
	  all the `^.aker$` words.
* [x] connect puzzler to solver
    * [x] gather statistics on iterations
	* [ ] There's a resource leak somewhere in `wordler/main` such that it slows
	  as it iterates when running with local dictionary. Find and fix that!
* [ ] problems: `/usr/share/dict` on Mac does not include...
    * [ ] plurals
	* [ ] fewer
	* [ ] heist
* [ ] load dictionary cross-platform
