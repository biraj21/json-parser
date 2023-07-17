package main

func CompareRuneSlices(r1 []rune, r2 []rune, n int) bool {
	if n > len(r1) || n > len(r2) {
		return false
	}

	for i := 0; i < n; i++ {
		if r1[i] != r2[i] {
			return false
		}
	}

	return true
}
