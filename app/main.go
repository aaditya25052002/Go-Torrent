package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"

	bencode "github.com/go-projects/go-torrent/bencode"
)

var _ = json.Marshal

func main() {
	data, err := os.ReadFile("sample.torrent")

	if err == nil {
		infoBytes, err := bencode.ExtractInfoBytes(data)
		infoBytesHash := sha1.Sum(infoBytes)

		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(infoBytesHash)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("error reading file: ", err)
		os.Exit(1)
	}
}
