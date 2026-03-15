package peers

import (
	"fmt"
	"net"
)

func WaitForUnChoke(conn net.Conn) {
	if _, err := waitForMessage(conn, MsgBitfield); err != nil {
		fmt.Println("error receiving bitfield: ", err)
		return
	}
	fmt.Println("Received bitfield")

	if err := sendMessage(conn, MsgInterested, nil); err != nil {
		fmt.Println("error sending interested: ", err)
		return
	}
	fmt.Println("Sent interested")

	if _, err := waitForMessage(conn, MsgUnchoke); err != nil {
		fmt.Println("error receiving unchoke: ", err)
		return
	}
	fmt.Println("Received unchoke")
}
