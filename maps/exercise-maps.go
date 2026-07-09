/*
	Zaimplementuj WordCount. Program powinien zwrócić mapę która policzy ilość wystąpień każdego "wyrazu" w stringu s.
 	Funkcja wc.Test uruchamia test twojej funkcji i wypisuje na ekranie "success" (pl. sukces) lub "failure" (pl. porażka).

	W tym ćwiczeniu może okazać się pomocne zapoznanie się z funkcją strings.Fields.
*/

package maps

import (
	"strings"
)

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	counts := make(map[string]int)

	for _, word := range words {
		counts[word]++
	}

	return counts
}
