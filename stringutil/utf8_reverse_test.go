package stringutil

import "testing"

func TestUtf8Reverse(t *testing.T) {
	for _, c := range []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{"Hello, 世界", "界世 ,olleH"},
		{"", ""},
	} {
		got := Utf8Reverse(c.in)
		if got != c.want {
			t.Errorf("Utf8Reverse(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func benchmarkUtf8Reverse(s string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		Utf8Reverse(s)
	}
}

func BenchmarkUtf8ReverseBytes(b *testing.B) {
	benchmarkUtf8Reverse("Hello, world", b)
}

func BenchmarkUtf8ReverseEmptyString(b *testing.B) {
	benchmarkUtf8Reverse("", b)
}

func BenchmarkUtf8ReverseRunes(b *testing.B) {
	benchmarkUtf8Reverse("Hello, 世界", b)
}
