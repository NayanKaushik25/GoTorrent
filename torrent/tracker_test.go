package torrent

import "testing"

func TestEscapeBinary(t *testing.T) {
	input := []byte{0x00, 0x20, 0x41, 0x7F, 0x80, 0xFF}
	got := escapeBinary(input)
	want := "%00%20A%7F%80%FF"
	if got != want {
		t.Fatalf("escapeBinary() = %q,want %q", got, want)
	}
}
