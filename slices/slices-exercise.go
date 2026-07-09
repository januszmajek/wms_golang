/*
	Zaimplementuj funkcję Pic. Powinna ona zwrócić wycinek o długości dy, którego elementami są wycinki
	zawierające dx 8-bitowych liczb całkowitych bez znaku (ang. 8-bit unisgned integer).
	Gdy uruchomisz program, wyświetli on obrazek, interpretując wartości intów w skali szarości
 	(tak naprawdę, to w skali „niebieskości”).

	Wybór obrazka należy do ciebie. Przykładowe funkcje które dają ciekawy rezultat to: (x+y)/2, x*y, and x^y.

	(Musisz użyć pętli by przypisać każdemu elementowi w [][]uint8 wartość typu []uint8.)

	(Użyj uint8(intValue) by dokonać konwersji typów.)
*/

package slices

func Pic(dx, dy int) [][]uint8 {
	var result [][]uint8
	for y := 0; y < dy; y++ {
		result = append(result, make([]uint8, dx))
		for x := 0; x < dx; x++ {
			result[y][x] = uint8((y + x) / 2)
		}
	}
	return result
}
