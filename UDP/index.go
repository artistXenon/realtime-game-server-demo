package UDP

import (
	"encoding/binary"
	"fmt"
	"inagame/crypto"
	"inagame/db"
	"inagame/state"
	"inagame/state/lobby"
	"math/big"
	"net"
	"time"
)

type Message struct {
	EventName string
	UserId    string
	Hash      string
	Body      string
}

type Header struct {
	Command byte
	User    *lobby.Player
	Lobby   *lobby.Lobby
}

// is this even reference or value
/*
 * hash(all of the following + online secret) 4
 * command 1
 * uid 10
 * lobbyid 5
 * [ body ~ ]
 */
func parseHeader(buf []byte) (header *Header, body *[]byte) {
	idInt := new(big.Int)
	idInt.SetBytes(buf[5:15])
	idString := idInt.String()
	lobbyString := string(buf[15:20])

	// TODO: if user not in memory, access database to fetch sessionkey and etc
	// note: if lobby where user is addresed to does not exist in this application, then we can reject connection (assume user does not exist)
	clientUser := lobby.Players[idString]
	clientLobby := state.Games[lobbyString]

	if clientLobby == nil { // shouldn't happen
		return nil, nil
	}

	// new player on this server. create and assign
	if clientUser == nil {
		sessionKey := db.GetPlayer(idString, lobbyString)
		if sessionKey == nil {
			return nil, nil
		}
		clientUser = lobby.CreatePlayer(idString, *sessionKey, clientLobby)

		// previous player joined new lobby. re assign player w/ refreshed session key
	} else if clientUser.Lobby.Id != clientLobby.Id {
		sessionKey := db.GetPlayer(idString, lobbyString)
		if sessionKey == nil {
			return nil, nil
		}
		clientUser.SessionKey = *sessionKey
		clientLobby.AssignPlayer(clientUser)
	}

	h := Header{Command: buf[4], User: clientUser}

	signedBody := append([]byte(clientUser.SessionKey), buf[4:]...)

	generatedHash := crypto.GenerateCRCHash(signedBody)
	// insert client session key for hash
	if generatedHash != binary.BigEndian.Uint32(buf[0:4]) {
		return nil, nil
	}

	h.Lobby = clientLobby

	bd := buf[20:]
	return &h, &bd
}

func UDPHandler(udpServer net.PacketConn, addr net.Addr, buf []byte, byteLength int) {
	header, body := parseHeader(buf)

	if header == nil {
		return
	}

	var response *[]byte
	doResponse := false
	var err error
	switch header.Command {
	case COMMAND_CONNECT:
		// do something
	case COMMAND_JOIN:
		err, response, doResponse = onJoin(header, body)
	case COMMAND_PING:
		err, response, doResponse = onPing(header, body)
	case COMMAND_PONG:
		err, _, doResponse = onPong(header, body)
	default:
		// something is wrong. ignore packet
	}

	if err != nil {
		fmt.Printf("[%v] %v", time.Now().Format(time.ANSIC), err) // TODO: print timestamp to
	}

	if doResponse {
		// TODO: maybe we need hashed response. to be considered far later
		udpServer.WriteTo(*response, addr)
	}
}
