package a

import (
	"os"
	"testing"
)

func T(t *testing.T, name string) (string, error) {
	t.Helper()
	return os.MkdirTemp("", name)
}

func B(b *testing.B, name string) (string, error) {
	b.Helper()
	return os.MkdirTemp("", name)
}

func TB(tb testing.TB, name string) (string, error) {
	tb.Helper()
	return os.MkdirTemp("", name)
}

func TDD(t *testing.T, name string) {
	t.Helper()
	for _, tt := range []struct {
		name string
	}{
		{name: "a"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			os.MkdirTemp("", tt.name)
		})
	}
}
