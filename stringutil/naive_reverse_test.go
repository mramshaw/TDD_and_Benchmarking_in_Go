package stringutil

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestNaiveReverse(t *testing.T) {
	for _, c := range []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{"Hello, 世界", "界世 ,olleH"},
		{"", ""},
	} {
		got := NaiveReverse(c.in)
		assert.Equalf(t, got, c.want, "NaiveReverse(%q) returned %q, wanted %q", c.in, got, c.want)
	}
}

func benchmarkNaiveReverse(s string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		NaiveReverse(s)
	}
}

func BenchmarkNaiveReverseBytes(b *testing.B) {
	benchmarkNaiveReverse("Hello, world", b)
}

func BenchmarkNaiveReverseEmptyString(b *testing.B) {
	benchmarkNaiveReverse("", b)
}

func BenchmarkNaiveReverseRunes(b *testing.B) {
	benchmarkNaiveReverse("Hello, 世界", b)
}
