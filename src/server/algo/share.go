// SDR - Labo 4
// Nicolas Crausaz & Maxime Scharwath
// Common functions used by multiple algorithms

package algo

import (
	"strings"
)

type Data = int

func CountLetter(word string, letter string) int {
	return strings.Count(
		strings.ToLower(word),
		strings.ToLower(letter),
	)
}
