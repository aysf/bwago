package main

import "testing"

func TestRun(t *testing.T) {
	err := Run()
	if err != nil {
		t.Error("fail Run()")
	}
}
