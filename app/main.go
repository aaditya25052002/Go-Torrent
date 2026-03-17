package main

import (
	"fmt"
	"os"

	"github.com/go-projects/go-torrent/torrent"
)

func main() {
	data, err := os.ReadFile("sample.torrent")

	if err == nil {
		outputBuffer, outputFileName, err := torrent.Run(data)
		if err != nil {
			fmt.Println("error running torrent: ", err)
			os.Exit(1)
		}
		outFile, _ := os.Create(fmt.Sprintf("files/%s", outputFileName))
		outFile.Write(outputBuffer)
	} else {
		fmt.Println("error reading file: ", err)
		os.Exit(1)
	}
}
