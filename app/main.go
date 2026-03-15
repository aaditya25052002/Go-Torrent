package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"

	bencode "github.com/go-projects/go-torrent/bencode"
	peers "github.com/go-projects/go-torrent/peers"
)

var _ = json.Marshal

func splitPieceHashes(pieces []byte) [][]byte {
	var hashes [][]byte

	for i := 0; i < len(pieces); i += 20 {
		hashes = append(hashes, pieces[i:i+20])
	}

	return hashes
}

func main() {
	data, err := os.ReadFile("sample.torrent")

	if err == nil {
		infoBytes, err := bencode.ExtractInfoBytes(data)
		infoBytesHash := sha1.Sum(infoBytes)

		if err != nil {
			fmt.Println(err)
			return
		}

		decoded, _, err := bencode.Decode(string(data))
		if err != nil {
			fmt.Println(err)
			return
		}

		torrent := decoded.(map[string]any)
		info := torrent["info"].(map[string]any)

		// pieces := []byte(info["pieces"].(string))
		// pieceHashes := splitPieceHashes(pieces)

		peersApiResponse := peers.Discover(torrent["announce"].(string), infoBytesHash[:], info["piece length"].(int))
		peersBinary := peersApiResponse["peers"].(string)
		peersList := peers.ParseBinary([]byte(peersBinary))

		fmt.Println("\nPeers Discovered: ")
		for _, peer := range peersList {
			fmt.Println(peer)
		}

		conn, err := peers.Connect(peersList[1], infoBytesHash[:], []byte("-GO0001-123456789012"))
		if err != nil {
			fmt.Println("error connecting to peer: ", err)
			return
		}
		defer conn.Close()

		pieceLength := info["piece length"].(int)
		pieces := []byte(info["pieces"].(string))
		pieceHashes := splitPieceHashes(pieces)

		pieceIndex := 0
		data, err := peers.DownloadPiece(conn, pieceIndex, pieceLength, pieceHashes[pieceIndex])
		if err != nil {
			fmt.Println("error downloading piece: ", err)
			return
		}
		os.WriteFile(fmt.Sprintf("piece_%d.bin", pieceIndex), data, 0644)
		fmt.Printf("Downloaded piece %d (%d bytes)\n", pieceIndex, len(data))
	} else {
		fmt.Println("error reading file: ", err)
		os.Exit(1)
	}
}
