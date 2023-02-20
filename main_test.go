package main

import (
	"strconv"
	"testing"
)

func TestGetPods(t *testing.T) {
	var results []string
	var err error

	results, err = GetPods("/Users/adam/code/go-getting-started/okteto-kube.config", "ascarb")
	if err != nil {
		t.Errorf("Error was not nil.")
	}

	if len(results) == 0 {
		t.Errorf("Results array not populated.")
	} else {
		t.Logf("Found pods")
	}
	t.Logf("Pods: " + strconv.Itoa(len(results)))
}
