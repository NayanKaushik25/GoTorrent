package torrent

import "fmt"

func findInfoBytes(data []byte) ([]byte, error) {
	pos := 0
	if pos >= len(data) || data[pos] != 'd' {
		return nil, fmt.Errorf("torrent must be a bencoded dictionary")
	}
	pos++
	for pos < len(data) {
		if data[pos] == 'e' {
			return nil, fmt.Errorf("torrent missing info dictionary")
		}
		key, next, err := scanString(data, pos)
		if err != nil {
			return nil, err
		}
		pos = next
		valueStart := pos
		valueEnd, err := scanValue(data, pos)
		if err != nil {
			return nil, err
		}

		if key == "info" {
			if data[valueStart] != 'd' {
				return nil, fmt.Errorf("torrent info must be a dictionary")
			}
			return data[valueStart:valueEnd], nil
		}

		pos = valueEnd
	}

	return nil, fmt.Errorf("unexpected end of torrent data")
}

func scanString(data []byte, pos int) (string, int, error) {
	start := pos

	for pos < len(data) && data[pos] != ':' {
		if data[pos] < '0' || data[pos] > '9' {
			return "", 0, fmt.Errorf("invalid bencoded string length")
		}
		pos++
	}

	if pos >= len(data) {
		return "", 0, fmt.Errorf("unterminated string length")
	}

	length := 0
	for i := start; i < pos; i++ {
		length = length*10 + int(data[i]-'0')
	}
	pos++ // consume ':'

	end := pos + length
	if end > len(data) {
		return "", 0, fmt.Errorf("truncated bencoded string")
	}

	return string(data[pos:end]), end, nil
}

func scanValue(data []byte, pos int) (int, error) {
	if pos >= len(data) {
		return 0, fmt.Errorf("unexpected end of bencoded value")
	}

	switch data[pos] {
	case 'i':
		pos++
		for pos < len(data) {
			if data[pos] == 'e' {
				return pos + 1, nil
			}
			pos++
		}
		return 0, fmt.Errorf("unterminated integer")

	case 'l', 'd':
		container := data[pos]
		pos++

		for pos < len(data) {
			if data[pos] == 'e' {
				return pos + 1, nil
			}

			if container == 'd' {
				_, next, err := scanString(data, pos)
				if err != nil {
					return 0, err
				}
				pos = next
			}

			next, err := scanValue(data, pos)
			if err != nil {
				return 0, err
			}
			pos = next
		}
		return 0, fmt.Errorf("unterminated container")

	default:
		_, end, err := scanString(data, pos)
		return end, err
	}
}
