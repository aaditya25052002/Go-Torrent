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

		totalLength := info["length"].(int)
		pieceLength := info["piece length"].(int)
		pieces := []byte(info["pieces"].(string))
		pieceHashes := splitPieceHashes(pieces)

		outputFileName := info["name"].(string)
		outFile, _ := os.Create(outputFileName)
		defer outFile.Close()

		peers.WaitForUnChoke(conn)

		outputBuffer := make([]byte, 0, totalLength)
		for i := 0; i < len(pieceHashes); i++ {
			thisLength := pieceLength
			if i == len(pieceHashes)-1 {
				thisLength = totalLength - (i * pieceLength)
			}
			data, err := peers.DownloadPiece(conn, i, thisLength, pieceHashes[i])
			if err != nil {
				fmt.Println("error downloading piece: ", err)
				return
			}
			outputBuffer = append(outputBuffer, data...)
			fmt.Printf("Downloaded piece %d (%d bytes)\n", i, len(data))
		}
		outFile.Write(outputBuffer)
	} else {
		fmt.Println("error reading file: ", err)
		os.Exit(1)
	}
}
