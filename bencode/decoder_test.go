package main

import (
	"bytes"
	"testing"
)

func TestDecodeInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"i42e", 42},
		{"i0e", 0},
		{"i-42e", -42},
	}
	i := 0
	for i < len(tests) {
		reader := bytes.NewReader([]byte(tests[i].input))
		got, err := Decode(reader)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got != tests[i].expected {
			t.Errorf("expected %v, got %v", tests[i].expected, got)
		}
		i++
	}
}

func TestDecodeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"4:spam", "spam"},
		{"3:abc", "abc"},
		{"10:helloworld", "helloworld"},
	}
	for _, test := range tests {
		reader := bytes.NewReader([]byte(test.input))
		got, err := Decode(reader)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != test.expected {
			t.Errorf("expected %q,got %q", test.expected, got)
		}
	}
}
