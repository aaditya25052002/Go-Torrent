package peers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	bencode "github.com/go-projects/go-torrent/bencode"
)

func Discover(announce string, infoHash []byte, pieceLength int) map[string]any {
	peerId := generatePeerId()

	params := url.Values{
		"info_hash":  {string(infoHash)},
		"peer_id":    {string(peerId)},
		"port":       {"6881"},
		"uploaded":   {"0"},
		"downloaded": {"0"},
		"left":       {strconv.Itoa(pieceLength)},
		"compact":    {"1"},
	}

	trackerUrl := announce + "?" + params.Encode()

	fmt.Println("tracking url: ", trackerUrl)

	resp, err := http.Get(trackerUrl)
	if err != nil {
		fmt.Println("error getting peers: ", err)
		return nil
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading body: ", err)
		return nil
	}

	decoded, _, _ := bencode.Decode(string(body))
	trackerResponse := decoded.(map[string]any)

	return trackerResponse
}

func generatePeerId() string {
	return "-GO0001-123456789012"
}
