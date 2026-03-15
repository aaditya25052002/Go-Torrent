package peers

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	MsgChoke         = 0
	MsgUnchoke       = 1
	MsgInterested    = 2
	MsgNotInterested = 3
	MsgHave          = 4
	MsgBitfield      = 5
	MsgRequest       = 6
	MsgPiece         = 7
	MsgCancel        = 8

	BlockSize = 16384 // 16 KiB — standard BitTorrent block size
)

type PeerMessage struct {
	ID      byte
	Payload []byte
}

func DownloadPiece(conn net.Conn, pieceIndex int, pieceLength int, pieceHash []byte) ([]byte, error) {
	requestBitfield(conn)
	requestInterested(conn)
	requestUnchoke(conn)
	requestBlocks(conn, pieceIndex, pieceLength)

	totalBlocks := (pieceLength + BlockSize - 1) / BlockSize
	pieceData := receivePieceData(conn, totalBlocks, pieceLength, pieceIndex)
	if pieceData == nil {
		return nil, fmt.Errorf("error receiving piece data")
	}

	if !isPieceIntegrityValid(pieceData, pieceHash) {
		return nil, fmt.Errorf("piece %d hash mismatch", pieceIndex)
	}

	return pieceData, nil
}

func requestBitfield(conn net.Conn) {
	if _, err := waitForMessage(conn, MsgBitfield); err != nil {
		fmt.Println("error receiving bitfield: ", err)
		return
	}
	fmt.Println("Received bitfield")
}

func requestInterested(conn net.Conn) {
	if err := sendMessage(conn, MsgInterested, nil); err != nil {
		fmt.Println("error sending interested: ", err)
		return
	}
	fmt.Println("Sent interested")
}

func requestUnchoke(conn net.Conn) {
	if _, err := waitForMessage(conn, MsgUnchoke); err != nil {
		fmt.Println("error receiving unchoke: ", err)
		return
	}
	fmt.Println("Received unchoke")
}

func requestBlocks(conn net.Conn, pieceIndex int, pieceLength int) {
	for offset := 0; offset < pieceLength; offset += BlockSize {
		blockLen := BlockSize
		if offset+blockLen > pieceLength {
			blockLen = pieceLength - offset // last block may be smaller
		}
		if err := requestBlock(conn, pieceIndex, offset, blockLen); err != nil {
			fmt.Println("error requesting block: ", err)
			return
		}
	}
}

func requestBlock(conn net.Conn, pieceIndex, offset, length int) error {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(pieceIndex))
	binary.BigEndian.PutUint32(payload[4:8], uint32(offset))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return sendMessage(conn, MsgRequest, payload)
}

func receivePieceData(conn net.Conn, totalBlocks int, pieceLength int, pieceIndex int) []byte {
	pieceData := make([]byte, pieceLength)

	received := 0
	for received < totalBlocks {
		msg, err := readMessage(conn)
		if err != nil {
			fmt.Println("error reading piece block: ", err)
			return nil
		}
		if msg == nil || msg.ID != MsgPiece {
			continue
		}

		idx := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
		begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
		block := msg.Payload[8:]

		if idx != pieceIndex {
			continue
		}

		copy(pieceData[begin:], block)
		received++
		fmt.Printf("  block %d/%d  offset=%d  size=%d\n", received, totalBlocks, begin, len(block))
	}
	return pieceData
}

func isPieceIntegrityValid(pieceData []byte, pieceHash []byte) bool {
	hash := sha1.Sum(pieceData)
	if !bytes.Equal(hash[:], pieceHash) {
		return false
	}
	return true
}

func readMessage(conn net.Conn) (*PeerMessage, error) {
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, lengthBuf); err != nil {
		return nil, fmt.Errorf("read msg length: %w", err)
	}

	length := binary.BigEndian.Uint32(lengthBuf)
	if length == 0 {
		return nil, nil // keep-alive
	}

	msgBuf := make([]byte, length)
	if _, err := io.ReadFull(conn, msgBuf); err != nil {
		return nil, fmt.Errorf("read msg body (%d bytes): %w", length, err)
	}

	return &PeerMessage{
		ID:      msgBuf[0],
		Payload: msgBuf[1:],
	}, nil
}

func sendMessage(conn net.Conn, id byte, payload []byte) error {
	length := uint32(1 + len(payload))
	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = id
	copy(buf[5:], payload)
	_, err := conn.Write(buf)
	return err
}

func waitForMessage(conn net.Conn, expectedID byte) (*PeerMessage, error) {
	for {
		msg, err := readMessage(conn)
		if err != nil {
			return nil, err
		}
		if msg == nil {
			continue // keep-alive
		}
		if msg.ID == expectedID {
			return msg, nil
		}
	}
}
