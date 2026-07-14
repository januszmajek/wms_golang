/*
Package maps
Zaimplementuj WordCount. Program powinien zwrócić mapę która policzy ilość wystąpień każdego "wyrazu" w stringu s.

W tym ćwiczeniu może okazać się pomocne zapoznanie się z funkcją strings.Fields.
*/
package maps

import (
	"strings"
)

// WordCount I have to proide comment i guess
func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	counts := make(map[string]int)

	for _, word := range words {
		counts[word]++
	}

	return counts
}
