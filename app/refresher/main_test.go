package main

import "testing"

var tests = []struct {
	name     string
	dividend float32
	diviser  float32
	result   float32
	hasError bool
}{
	{"valid-data", 100, 10, 10, false},
	{"invalid-data", 100, 0, 0, true},
	{"expected-5", 50, 10, 5, false},
}

func TestDivision(t *testing.T) {
	for _, tt := range tests {
		result, err := divide(tt.dividend, tt.diviser)
		if tt.hasError && err == nil {
			t.Error("Got an unexpected result")
		} else if !tt.hasError && err != nil {
			t.Error("Got an unexpected result")
		}

		if result != tt.result {
			t.Error("Expected results not match")
		}
	}
}

func TestDivide(t *testing.T) {
	_, err := divide(10, 10)
	if err != nil {
		t.Error("Got an error we should not have")
	}
}

func TestHardDivide(t *testing.T) {
	_, err := divide(10, 0)
	if err == nil {
		t.Error("Did not have an error we should have")
	}
}
