package parsers

import (
	"testing"
)

func Test_Changes(t *testing.T) {
	changes := []string{"hello/myworld", "whatever"}
	res := ChangedPaths("whatever,hello/myworld", changes)
	isTrue(t, res["hello/myworld"])
	isTrue(t, res["whatever"])
	isTrue(t, !res["something"])
}

func isTrue(t *testing.T, b bool) {
	t.Helper()
	if b != true {
		t.Fatal("Expected true but was false")
	}
}
