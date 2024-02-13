package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// initialize a new time.Time object and pass it to the humanDate function.
	tm := time.Date(2024, 2, 13, 12, 20, 0, 0, time.UTC)
	hd := humanDate(tm)

	// Check that the output from humanDate function is in the format
	// we expect. If it isn't what we expect, we use the t.Errorf() function to
	// indicate that the test has failed and log the expected and actual values
	if hd != "13 Feb 2024 at 12:20 UTC" {
		t.Errorf("got %q, want %q", hd, "13 Feb 2024 at 12:20 UTC") // use of "quote" formatting verb
	}

}
