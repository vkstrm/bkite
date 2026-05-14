package main

import "testing"

func TestAdd(t *testing.T) {
	expected := 2
	actual := add(1, 1)
	if expected != actual {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}
