package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/NayanKaushik25/GoTorrent/bencode"
)

type MetaInfo struct {
	Announce    string
	Name        string
	PieceLength int64
	Pieces      string
	Length      int64
	InfoHash    [20]byte
}

func ParseMetaInfo(r io.Reader) (*MetaInfo, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	decoded, err := bencode.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	dict, ok := decoded.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("torrent file must contain a dictionary")
	}
	infoValue, ok := dict["info"]
	if !ok {
		return nil, fmt.Errorf("Torrent missing \"info\" dictionary")
	}
	info, ok := infoValue.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("torrent info must be in a dictionary")
	}
	announceVal, ok := dict["announce"]
	if !ok {
		return nil, fmt.Errorf("torrent missing announce URL")
	}
	announce, ok := announceVal.(string)
	if !ok {
		return nil, fmt.Errorf("announce must be a string")
	}

	// Extract name
	nameValue, ok := info["name"]
	if !ok {
		return nil, fmt.Errorf("torrent missing name")
	}
	name, ok := nameValue.(string)
	if !ok {
		return nil, fmt.Errorf("name must be a string")
	}

	// Extract piece length
	pieceLengthValue, ok := info["piece length"]
	if !ok {
		return nil, fmt.Errorf("torrent missing piece length")
	}
	pieceLength, ok := pieceLengthValue.(int64)
	if !ok {
		return nil, fmt.Errorf("piece length must be an integer")
	}

	// Extract pieces
	piecesValue, ok := info["pieces"]
	if !ok {
		return nil, fmt.Errorf("torrent missing pieces")
	}
	pieces, ok := piecesValue.(string)
	if !ok {
		return nil, fmt.Errorf("pieces must be a string")
	}

	// Extract length
	lengthValue, ok := info["length"]
	if !ok {
		return nil, fmt.Errorf("torrent missing length")
	}
	length, ok := lengthValue.(int64)
	if !ok {
		return nil, fmt.Errorf("length must be an integer")
	}

	infoBytes, err := findInfoBytes(data)
	if err != nil {
		return nil, err
	}

	infoHash := sha1.Sum(infoBytes)

	return &MetaInfo{
		Announce:    announce,
		Name:        name,
		PieceLength: pieceLength,
		Pieces:      pieces,
		Length:      length,
		InfoHash:    infoHash,
	}, nil
}
