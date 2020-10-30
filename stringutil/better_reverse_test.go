package stringutil

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestBetterReverse(t *testing.T) {
	for _, c := range []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{"Hello, 世界", "界世 ,olleH"},
		{"", ""},
	} {
		got := BetterReverse(c.in)
		assert.Equalf(t, got, c.want, "BetterReverse(%q) returned %q, wanted %q", c.in, got, c.want)
	}
}

func benchmarkBetterReverse(s string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		BetterReverse(s)
	}
}

func BenchmarkBetterReverseBytes(b *testing.B) {
	benchmarkBetterReverse("Hello, world", b)
}

func BenchmarkBetterReverseEmptyString(b *testing.B) {
	benchmarkBetterReverse("", b)
}

func BenchmarkBetterReverseRunes(b *testing.B) {
	benchmarkBetterReverse("Hello, 世界", b)
}
