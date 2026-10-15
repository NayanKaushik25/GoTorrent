package bencode

import "io"

type CountingReader struct {
	Reader io.Reader
	Count  int
}

func (cr *CountingReader) Read(p []byte) (int, error) {
	n, err := cr.Reader.Read(p)
	cr.Count += n
	return n, err
}
