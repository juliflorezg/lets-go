package assert

import "testing"

// generic function, we’ll be able to use it irrespective
// of what the type of the 'actual' and 'expected' values is
func Equal[T comparable](t *testing.T, actual, expected T) {

	t.Helper()

	if actual != expected {
		t.Errorf("got %v, want %v", actual, expected)
	}
}
