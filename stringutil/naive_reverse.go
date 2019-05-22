package stringutil

// NaiveReverse returns its argument string reversed rune-wise left to right.
func NaiveReverse(s string) string {
	r := []rune(s)
    n := len(r)
	rOut := make([]rune, n)
	for i := 0; i < n; i += 1 {
		rOut[i] = r[n - 1 - i]
	}
	return string(rOut)
}
