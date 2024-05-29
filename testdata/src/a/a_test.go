package a_test

import (
	"os"
	"testing"
)

func TestT(t *testing.T) {
	t.Helper()
	os.MkdirTemp("", "test") // want "os\\.MkdirTemp\\(\\) can be replaced by `t.TempDir\\(\\)`"
}

func BenchmarkB(b *testing.B) {
	b.Helper()
	os.MkdirTemp("", "test") // want "os\\.MkdirTemp\\(\\) can be replaced by `b.TempDir\\(\\)`"
}

func TestTDD(t *testing.T) {
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
