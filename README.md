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
	* Use weighting, but use wordler's `hard mode`, which is actually easier
	  than my initial implementation.
	    * 89.0% success rate
		* Average guesses to win: 4.34 guesses

## TODO
* [ ] Optimizations.
    * [ ] `egrep ^.ater$ /usr/share/dict/words` reveals a
      pathological case where `bater` can take 8 guesses. Optimized solver
      should seek to avoid these pathological cases by avoiding guesses that
      use a pattern of letters repeated in other words. Another test case:
      `baker` is solved ~34% of the time with the current algorithm due to
      all the `^.aker$` words.
      NOTE: this TODO is based on `dict`, not the wordle dictionary.
    * [ ] Running iterator with `--iterations=1000 --solution=glass` only wins
	  48.1% of games played.
* [x] Easier `Hard Mode`. It turns out that my `Hard Mode` is significantly
      harder than Wordler's, and following theirs allows additional optimizations.
      For example, if guess `arose` results in `a` and `o` identified as in
      the puzzle but in the wrong place, Wordler alows guess `atoms` to follow.
      I don't. The optimization would be for solver to maintain a list of
      possible solutions (my `Hard Mode`) separate from a list of permissible
      guesses (Wordler's `Hard Mode`) and then identify the optimal guess based
      on the possible solution set.
	  `Hard Mode` requires merely that:
	  1. Letters known to be in the right place must stay there.
	  2. Letters known to be in the puzzle must be included in your guess, but
	     may continue to be in the place known to be incorrect.
* [ ] Fix resource leak in the platform dictionary loader -- iterator slows 
      as it iterates when running with platform dictionary.
* [ ] load dictionary cross-platform
