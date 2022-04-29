package wordler

import (
	"bufio"
	"os"
)

type Dictionary struct {
	*WordList
}

// NewDictionary returns a Dictionary containing /usr/share/dict/words.
// TODO: make this platform-independent via goos.Is*
func NewDictionary() (Dictionary, error) {
	file, err := os.Open(`/usr/share/dict/words`)
	if err != nil {
		return Dictionary{}, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return Dictionary{NewWordList(lines)}, scanner.Err()
}
