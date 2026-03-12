package main

import (
	"encoding/json"
	"fmt"
	"os"

	bencode "github.com/go-projects/go-torrent/bencode"
)

var _ = json.Marshal

func main() {
	data, err := os.ReadFile("sample.torrent")

	if err == nil {
		bencodedValue := string(data)

		decoded, _, err := bencode.Decode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("error reading file: ", err)
		os.Exit(1)
	}
}
