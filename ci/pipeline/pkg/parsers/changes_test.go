package parsers

import (
	"fmt"
	"testing"
)

func Test_Changes(t *testing.T) {
	changes := []string{"hello/myworld", "whatever"}
	res := ChangedPaths("whatever,hello/myworld", changes)
	fmt.Println(res)
	if res["hello/myworld"] != true {
		t.FailNow()
	}
	if res["whatever"] != true {
		t.FailNow()
	}
	if res["something"] != false {
		t.FailNow()
	}
}
