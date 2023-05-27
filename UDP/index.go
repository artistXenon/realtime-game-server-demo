package UDP

import (
	"net"
	"strconv"
	"strings"
)

type Message struct {
	EventName string
	UserId    int64
	Hash      string
	Body      string
}

func ParseMessage(rawMsg []byte) *Message {
	// create message
	msg := string(rawMsg)

	// search for message title
	idx := strings.Index(msg, "!")
	m := new(Message)
	if idx == -1 {
		m.EventName = "!"
	}
	m.EventName = msg[:idx]

	// search for sender id
	msg = msg[idx+1:]
	idx = strings.Index(msg, "!")
	if idx == -1 {
		m.UserId = -1
	} else {
		id, err := strconv.ParseInt(msg[:idx], 0, 64)
		if err != nil {
			m.UserId = -1
		}
		m.UserId = id
	}

	// search for hash associated
	// todo: this should be verified with the body
	msg = msg[idx+1:]
	idx = strings.Index(msg, "!")
	if idx == -1 {
		m.Hash = ""
	} else {
		m.Hash = msg[:idx]
	}

	//search for body
	m.Body = msg[idx+2:]
	return m
}

func UDPHandler(udpServer net.PacketConn, addr net.Addr, buf []byte, byteLength int) {
	m := ParseMessage(buf[0:byteLength])

	response := ""
	doResponse := false
	switch m.EventName {
	case "connect":
		// do something
	case "join":
		response, doResponse = onJoin(m)
		if response[0] != '!' {

		}
		// do something
	case "ping":
		response, doResponse = onPing(m)
	case "pong":
		onPong(m)
		doResponse = false
	default:
		//something is wrong. ignore packet
	}

	// time := time.Now().Format(time.ANSIC)
	// fmt.Printf("time received: %v\n", time)
	// fmt.Printf("msg received: %v\n", *m)

	// fmt.Println(len(string(buf)))

	if doResponse {
		udpServer.WriteTo([]byte(response), addr)
	}
}

// message format
/**
[evenetname]![userid]![hash]!
{data json}
*/
