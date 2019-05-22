package stringutil

// BetterReverse returns its argument string reversed rune-wise left to right.
func BetterReverse(s string) string {
	r := []rune(s)
    n := len(r)
	rOut := make([]rune, n)
	i := n - 1
    for _, c := range s {
		rOut[i] = c
        i -= 1
    }
	return string(rOut)
}
