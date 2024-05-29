package a

import (
	"os"
	"testing"
)

func T(t *testing.T, name string) (string, error) {
	t.Helper()
	os.MkdirTemp("", name)        // want "os\\.MkdirTemp\\(\\) can be replaced by `t\\.TempDir\\(\\)`"
	return os.MkdirTemp("", name) // want "os\\.MkdirTemp\\(\\) can be replaced by `testing\\.TempDir\\(\\)`"
}

func B(b *testing.B, name string) (string, error) {
	b.Helper()
	os.MkdirTemp("", name)        // want "os\\.MkdirTemp\\(\\) can be replaced by `b\\.TempDir\\(\\)`"
	return os.MkdirTemp("", name) // want "os\\.MkdirTemp\\(\\) can be replaced by `testing\\.TempDir\\(\\)`"
}

func TB(tb testing.TB, name string) (string, error) {
	tb.Helper()
	os.MkdirTemp("", name)        // want "os\\.MkdirTemp\\(\\) can be replaced by `tb.TempDir\\(\\)`"
	return os.MkdirTemp("", name) // want "os\\.MkdirTemp\\(\\) can be replaced by `testing.TempDir\\(\\)`"
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
			os.MkdirTemp("", tt.name) // want "os\\.MkdirTemp\\(\\) can be replaced by `t.TempDir\\(\\)` in anonymous function"
		})
	}
}
