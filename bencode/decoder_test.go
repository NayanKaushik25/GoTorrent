package main

import (
	"bytes"
	"testing"
)

func TestDecodeInteger(t *testing.T) {
	reader := bytes.NewReader([]byte("i42e"))
	got, err := Decode(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != int64(42) {
		t.Errorf("expected 42, got %v", got)
	}
}
