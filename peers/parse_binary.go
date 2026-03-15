package peers

import (
	"encoding/binary"
	"net"
)

type Peer struct {
	IP   string
	Port int
}

func ParseBinary(peersBinary []byte) []Peer {
	var peers []Peer

	for i := 0; i < len(peersBinary); i += 6 {
		ip := net.IP(peersBinary[i : i+4])
		port := binary.BigEndian.Uint16(peersBinary[i+4 : i+6])

		peers = append(peers, Peer{IP: ip.String(), Port: int(port)})
	}

	return peers
}
