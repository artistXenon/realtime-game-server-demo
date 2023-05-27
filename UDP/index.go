package UDP

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func UDPHandler(udpServer net.PacketConn, addr net.Addr, buf []byte, byteLength int) {
	decodedMsg := string(buf[0:byteLength])
	eventIdx := strings.Index(decodedMsg, "!")
	if eventIdx == -1 {
		// something's wrong
		return
	}
	eventName := decodedMsg[:eventIdx]
	eventBody := decodedMsg[eventIdx+1:]

	response := ""
	doResponse := false
	switch eventName {
	case "connect":
		// do something
	case "join":
		response, doResponse = onJoin(eventBody)
		if response[0] != '!' {

		}
		// do something
	case "ping":
		onPing(eventBody)
	default:
		//something is wrong. ignore packet
	}

	time := time.Now().Format(time.ANSIC)
	fmt.Printf("time received: %v\n", time)

	// fmt.Println(len(string(buf)))

	if doResponse {
		udpServer.WriteTo([]byte(response), addr)
	}
}

// message format
/**
[evenetname]!{data json}
*/
