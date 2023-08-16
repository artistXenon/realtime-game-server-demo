package TCP

import (
	"encoding/binary"
	"fmt"
	"inagame/state/lobby"
	"io"
	"log"
	"net"
	"time"
)

func TCPHandler(conn net.Conn) {
	var player *lobby.Player
	isPlayerBoundConnection := false

	buf := make([]byte, 1024)

	log.Printf("connection made %p", conn)

	for {
		length, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("connection closed %p", conn)
				// connection closed
			} else {
				log.Printf("unresolved error: %v", err)
				time.Sleep(5 * time.Second)
			}
			break
		}

		if length <= 0 {
			continue
		}

		var command byte

		if !isPlayerBoundConnection {
			command = buf[4]
			joinResult, errcode := joinHandler(&buf, length)
			var errMsg string
			switch errcode {
			case 0:
				player = joinResult
				if player.TCP != nil && player.TCP != &conn {
					(*player.TCP).Write([]byte{0x00})
					// (*player.TCP).Close()
				}
				player.TCP = &conn
				isPlayerBoundConnection = true
				fmt.Printf("join-tcp: %s joined lobby %s\n", player.Id, player.Lobby.Id)
			case 1:
				errMsg = "matching user/lobby not found"
			case 2:
				errMsg = "invalid hash on join request"
			}
			conn.Write([]byte{command, 0, 0, 0, 1, byte(errcode)})

			if errcode != 0 {
				log.Printf("tcp join has failed due to error: %s", errMsg)
				// invalid authentication from client.
				go func() {
					time.Sleep(time.Second * 1)
					conn.Close()
				}()
				return
			} else {
				continue
			}
		}
		command = buf[0]
		body := buf[1:length]
		var res *[]byte
		switch command {
		case COMMAND_LOBBY:
			res, err = lobbyHandler(&body, player) // todo" eat player/lobby
		}

		var resLength int32 = int32(len(*res))
		resMeta := make([]byte, 5)
		resMeta[0] = command
		binary.BigEndian.PutUint32(resMeta[1:], uint32(resLength))

		if err != nil {
			log.Printf("state error: %v", err)
		} else {
			conn.Write(append(resMeta, *res...))
		}

		// to some stuffs for communication.
		// conn.Write()
	}
}
