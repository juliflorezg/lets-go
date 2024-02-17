package assert

import (
	"strings"
	"testing"
)

// generic function, we’ll be able to use it irrespective
// of what the type of the 'actual' and 'expected' values is
func Equal[T comparable](t *testing.T, actual, expected T) {

	t.Helper()

	if actual != expected {
		t.Errorf("got %v, want %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()
	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got %v, want %v", actual, expectedSubstring)
	}
}
