package torrent

import (
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
}

func ParseMetaInfo(r io.Reader) (*MetaInfo, error) {
	decoded, err := bencode.Decode(r)
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

	return &MetaInfo{
		Announce:    announce,
		Name:        name,
		PieceLength: pieceLength,
		Pieces:      pieces,
		Length:      length,
	}, nil
}
