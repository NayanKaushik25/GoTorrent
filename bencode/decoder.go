package main

import (
	"fmt"
	"io"
	"strconv"
)

func decodeInteger(r io.Reader) (int64, error) {
	decoded := ""
	for {
		buf := make([]byte, 1)
		n, err := r.Read(buf)
		if err != nil {
			return 0, err
		}
		if n == 1 {
			if buf[0] != 'e' {
				decoded += string(buf[0])
			} else {
				break
			}
		} else {
			return 0, fmt.Errorf("failed to read a byte")
		}
	}
	return strconv.ParseInt(decoded, 10, 64)
}

func decodeList(r io.Reader) ([]any, error) {
	panic("not implemented")
}

func decodeDictionary(r io.Reader) (map[string]any, error) {
	panic("not implemented")
}

func decodeString(r io.Reader, firstByte byte) (string, error) {
	panic("not implemented")
}

func Decode(r io.Reader) (any, error) {
	buf := make([]byte, 1)
	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	if n == 1 {
		switch buf[0] {
		case 'i':
			return decodeInteger(r)
		case 'l':
			return decodeList(r)
		case 'd':
			return decodeDictionary(r)
		default:
			return decodeString(r, buf[0])
		}
	}
	return nil, fmt.Errorf("failed to read a byte")
}
