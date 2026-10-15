package torrent

import (
	"strings"
	"testing"
)

func TestParseMetaInfo(t *testing.T) {
	input := "d8:announce27:https://tracker.example.com4:infod6:lengthi5e4:name8:test.txt12:piece lengthi16384e6:pieces20:abcdefghijklmnopqrstee"

	got, err := ParseMetaInfo(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseMetaInfo() error = %v", err)
	}

	if got.Announce != "https://tracker.example.com" {
		t.Errorf("Announce = %q", got.Announce)
	}
	if got.Name != "test.txt" {
		t.Errorf("Name = %q", got.Name)
	}
	if got.PieceLength != 16384 {
		t.Errorf("PieceLength = %d", got.PieceLength)
	}
	if got.Length != 5 {
		t.Errorf("Length = %d", got.Length)
	}
	if got.Pieces != "abcdefghijklmnopqrst" {
		t.Errorf("Pieces = %q", got.Pieces)
	}
}
