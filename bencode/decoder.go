package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func readByte(r io.Reader) (byte, error) {
	buf := make([]byte, 1)

	n, err := r.Read(buf)
	if err != nil {
		return 0, err
	}

	if n != 1 {
		return 0, fmt.Errorf("failed to read a byte")
	}

	return buf[0], nil
}

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
	if len(decoded) > 1 && decoded[0] == '0' {
		return 0, fmt.Errorf("leading zeroes not permitted")
	}
	if strings.HasPrefix(decoded, "-0") {
		return 0, fmt.Errorf("invalid negative zero or leading zero")
	}
	return strconv.ParseInt(decoded, 10, 64)
}

func decodeList(r io.Reader) ([]any, error) {
	var result []any
	for {
		b, err := readByte(r)
		if err != nil {
			return nil, err
		}
		if b == 'e' {
			return result, nil
		}
		value, err := decodeValue(r, b)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
}

func decodeDictionary(r io.Reader) (map[string]any, error) {
	result := map[string]any{}
	for {
		b, err := readByte(r)
		if err != nil {
			return nil, err
		}
		if b == 'e' {
			return result, nil
		}
		if b < '0' || b > '9' {
			return nil, fmt.Errorf("format error")
		}
		key, err := decodeString(r, b)
		if err != nil {
			return nil, err
		}
		valueByte, err := readByte(r)
		if err != nil {
			return nil, err
		}
		value, err := decodeValue(r, valueByte)
		if err != nil {
			return nil, err
		}
		result[key] = value
	}
}

func decodeString(r io.Reader, firstByte byte) (string, error) {
	lengthStr := string(firstByte)
	for {
		b, err := readByte(r)
		if err != nil {
			return "", err
		}
		if b == ':' {
			break
		}
		lengthStr += string(b)
	}
	if len(lengthStr) > 1 && lengthStr[0] == '0' {
		return "", fmt.Errorf("leading zeros not permitted in string length")
	}
	for _, ch := range lengthStr {
		if ch < '0' || ch > '9' {
			return "", fmt.Errorf("string length must contain only digits")
		}
	}
	length, err := strconv.ParseInt(lengthStr, 10, 64)
	if err != nil {
		return "", err
	}
	if length < 0 {
		return "", fmt.Errorf("String cannot have negative length")
	}
	result := ""
	for i := int64(0); i < length; i++ {
		b, err := readByte(r)
		if err != nil {
			return "", err
		}
		result += string(b)
	}
	return result, nil
}

func decodeValue(r io.Reader, firstByte byte) (any, error) {
	switch firstByte {
	case 'i':
		return decodeInteger(r)
	case 'l':
		return decodeList(r)
	case 'd':
		return decodeDictionary(r)
	default:
		return decodeString(r, firstByte)
	}
}

func Decode(r io.Reader) (any, error) {
	firstByte, err := readByte(r)
	if err != nil {
		return nil, err
	}
	return decodeValue(r, firstByte)
}
