package algo

import "strings"

type Data = int

func CountLetter(word string, letter string) int {
	counter := 0
	letter = strings.ToLower(letter)
	for _, l := range strings.ToLower(word) {
		if string(l) == letter {
			counter++
		}
	}
	return counter
}
