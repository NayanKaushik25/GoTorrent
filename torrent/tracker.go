package torrent

import (
	"net/url"
	"strconv"
	"strings"
)

type Peer struct {
	IP   string
	Port uint16
}

type TrackerResponse struct {
	Interval int
	Peers    []Peer
}

func escapeBinary(data []byte) string {
	const hex = "0123456789ABCDEF"
	var result strings.Builder
	for _, b := range data {
		if (b >= 'a' && b <= 'z') ||
			(b >= 'A' && b <= 'Z') ||
			(b >= '0' && b <= '9') ||
			b == '-' || b == '_' || b == '.' || b == '~' {
			result.WriteByte(b)
		} else {
			result.WriteByte('%')
			result.WriteByte(hex[b>>4])
			result.WriteByte(hex[b&0x0F])
		}
	}
	return result.String()
}

func buildTrackerURL(announce string, infoHash [20]byte, peerID string, port uint16, uploaded, downloaded, left int64) (string, error) {
	u, err := url.Parse(announce)
	if err != nil {
		return "", err
	}
	params := u.Query()
	params.Set("peer_id", peerID)
	params.Set("port", strconv.Itoa(int(port)))
	params.Set("uploaded", strconv.FormatInt(uploaded, 10))
	params.Set("downloaded", strconv.FormatInt(downloaded, 10))
	params.Set("left", strconv.FormatInt(left, 10))
	params.Set("compact", "1")
	u.RawQuery = params.Encode()
	return u.String(), nil
}
