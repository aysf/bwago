package main

import "testing"

func TestMain(t *testing.T) {
	err := Run()
	if err != nil {
		t.Error("fail Run()")
	}
}
