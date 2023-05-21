package UDP

import (
	"fmt"
	"net"
	"time"
	"strings"
)

func UDPHandler(udpServer net.PacketConn, addr net.Addr, buf []byte, byteLength int) {
	decodedMsg := string(buf[0:byteLength])
    eventIdx := strings.Index(decodedMsg, "!")
	if eventIdx == -1 {
		// something's wrong
		return
	}
	eventName := decodedMsg[:eventIdx]
	eventBody := decodedMsg[eventIdx + 1:]

	response := ""
	doResponse := false
	switch eventName {
	// events that does not require body
		case "connect":
			// do something

	// events that require body
		case "join":
			response, doResponse = onJoin(eventBody)
			// do something
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