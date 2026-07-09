package fibonacci

/*
	Pobawmy się trochę funkcjami.
	Zaimplementuj funkcję fibonacci która zwróci funkcję (domknięcie) która zwraca ciąg liczby Fibonacciego (0, 1, 1, 2, 3, 5, ...).
*/

// fibonacci to funkcja która zwraca
// funkcję która zwraca int.
func Fibonacci() func() int {
	numbers := []int{0}

	return func() int {
		if len(numbers) < 2 {
			numbers = append(numbers, 1)
			return 0
		}
		if len(numbers) == 2 {
			numbers = append(numbers, 1)
			return 1
		}
		if len(numbers) == 3 {
			numbers = append(numbers, 1)
			return 1
		}
		next := numbers[len(numbers)-1] + numbers[len(numbers)-2]
		numbers = append(numbers, next)
		return next
	}
}

/*
	bardziej eleganckie
func fibonacci() func() int {
	a, b := 0, 1

	return func() int {
		result := a
		a, b = b, a+b
		return result
	}
}
*/
