package TCP

import (
	"fmt"
	"inagame/state/lobby"
	"io"
	"log"
	"net"
)

func TCPHandler(conn net.Conn) {
	var player *lobby.Player
	isPlayerBoundConnection := false

	buf := make([]byte, 1024)

	for {
		length, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("connection closed")
				// connection closed
			} else {
				log.Printf("unresolved error: %v", err)
			}
			break
		}

		if length <= 0 {
			continue
		}

		if !isPlayerBoundConnection {
			joinResult, err := joinHandler(&buf, length)
			if err != nil {
				log.Printf("tcp join has failed due to error: %v", err)
				// invalid authentication from client.
				conn.Close()
				return
			}
			player = joinResult
			player.TCP = &conn
			isPlayerBoundConnection = true
			fmt.Printf("%s joined lobby %s\n", player.Id, player.Lobby.Id)
			continue
		}
		command := buf[0]
		body := buf[1:length]
		var res *[]byte
		switch command {
		case COMMAND_STATE:
			res, err = stateHandler(&body)
		}

		if err != nil {
			log.Printf("state error: %v", err)
		} else {
			conn.Write(*res)
		}

		// to some stuffs for communication.
		// conn.Write()
	}

}
