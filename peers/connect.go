package peers

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

func buildHandshake(infoHash []byte, peerId []byte) []byte {
	buf := make([]byte, 68)
	buf[0] = 19
	copy(buf[1:20], "BitTorrent protocol")
	copy(buf[28:48], infoHash)
	copy(buf[48:68], peerId)
	return buf
}

func Connect(peer Peer, infoHash []byte, peerId []byte) (net.Conn, error) {
	addr := net.JoinHostPort(peer.IP, fmt.Sprintf("%d", peer.Port))

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", addr, err)
	}

	handshake := buildHandshake(infoHash, peerId)
	if _, err := conn.Write(handshake); err != nil {
		conn.Close()
		return nil, fmt.Errorf("send handshake: %w", err)
	}

	resp := make([]byte, 68)
	if _, err := io.ReadFull(conn, resp); err != nil {
		conn.Close()
		return nil, fmt.Errorf("read handshake: %w", err)
	}

	if !bytes.Equal(resp[28:48], infoHash) {
		conn.Close()
		return nil, fmt.Errorf("info hash mismatch from %s", addr)
	}

	fmt.Printf("Handshake OK with %s (peer id: %x)\n", addr, resp[48:68])
	return conn, nil
}
