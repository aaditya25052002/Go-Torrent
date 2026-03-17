package torrent

import (
	"crypto/sha1"
	"fmt"

	"github.com/go-projects/go-torrent/bencode"
	"github.com/go-projects/go-torrent/peers"
)

type PieceJob struct {
	index  int
	hash   []byte
	length int
}

type PieceResult struct {
	index int
	data  []byte
}

func Run(data []byte) ([]byte, string, error) {
	infoBytes, err := bencode.ExtractInfoBytes(data)
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}
	infoBytesHash := sha1.Sum(infoBytes)

	torrentInfo, err := decodeTorrentInfo(data)
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}
	peersList := disoverPeers(torrentInfo["announce"].(string), infoBytesHash[:], torrentInfo["piece_length"].(int))

	outputBuffer := downloadPieces(peersList, infoBytesHash[:], []byte(torrentInfo["pieces"].(string)), torrentInfo["piece_length"].(int), torrentInfo["total_length"].(int))

	return outputBuffer, torrentInfo["name"].(string), nil
}

func splitPieceHashes(pieces []byte) [][]byte {
	var hashes [][]byte

	for i := 0; i < len(pieces); i += 20 {
		hashes = append(hashes, pieces[i:i+20])
	}

	return hashes
}

func disoverPeers(announce string, infoBytesHash []byte, pieceLength int) []peers.Peer {
	peersApiResponse := peers.Discover(announce, infoBytesHash, pieceLength)
	peersBinary := peersApiResponse["peers"].(string)
	peersList := peers.ParseBinary([]byte(peersBinary))

	fmt.Println("\nPeers Discovered: ")
	for _, peer := range peersList {
		fmt.Println(peer)
	}

	return peersList
}

func downloadPieces(peersList []peers.Peer, infoBytesHash []byte, pieces []byte, pieceLength int, totalLength int) []byte {
	pieceHashes := splitPieceHashes(pieces)
	jobs := make(chan PieceJob, len(pieceHashes))
	results := make(chan PieceResult, len(pieceHashes))

	for i, hash := range pieceHashes {
		thisLength := pieceLength
		if i == len(pieceHashes)-1 {
			thisLength = totalLength - (i * pieceLength)
		}
		jobs <- PieceJob{index: i, hash: hash, length: thisLength}
	}

	for _, peer := range peersList {
		go func(peer peers.Peer) {
			conn, err := peers.Connect(peer, infoBytesHash, []byte("-GO0001-123456789012"))
			if err != nil {
				fmt.Println("error connecting to peer: ", err)
				return
			}
			defer conn.Close()

			peers.WaitForUnChoke(conn)

			for {
				job, ok := <-jobs
				if !ok {
					return
				}
				data, err := peers.DownloadPiece(conn, job.index, job.length, job.hash)
				if err != nil {
					fmt.Println("error downloading piece: ", err)
					jobs <- job
					continue
				}
				results <- PieceResult{index: job.index, data: data}
			}
		}(peer)
	}

	outputBuffer := make([]byte, totalLength)
	for i := 0; i < len(pieceHashes); i++ {
		result := <-results
		copy(outputBuffer[result.index*pieceLength:], result.data)
	}
	close(jobs)
	return outputBuffer
}

func decodeTorrentInfo(data []byte) (map[string]any, error) {
	decoded, _, err := bencode.Decode(string(data))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	root := decoded.(map[string]any)
	info := root["info"].(map[string]any)

	torrentInfo := map[string]any{
		"piece_length": info["piece length"].(int),
		"total_length": info["length"].(int),
		"name":         info["name"].(string),
		"pieces":       info["pieces"].(string),
		"announce":     root["announce"].(string),
	}
	return torrentInfo, nil
}
