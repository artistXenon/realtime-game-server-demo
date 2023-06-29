package UDP

import (
	"encoding/binary"
	"fmt"
	"inagame/crypto"
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
	Count   uint32
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
func parseHeader(buf *[]byte, byteLength int) (header *Header, body *[]byte) {
	packetCount := binary.BigEndian.Uint32((*buf)[0:4])
	idInt := new(big.Int)
	idInt.SetBytes((*buf)[9:19])
	idString := idInt.String()
	lobbyString := string((*buf)[19:24])

	clientUser := lobby.Players[idString]
	clientLobby := state.Games[lobbyString]

	// new player on this server. create and assign
	if clientUser == nil || clientLobby == nil || clientUser.Lobby.Id != clientLobby.Id {
		return nil, nil
		// previous player joined new lobby. re assign player w/ refreshed session key
	}

	h := Header{Count: packetCount, Command: (*buf)[8], User: clientUser}

	skipHash := h.Command == COMMAND_PING || h.Command == COMMAND_PONG

	if !skipHash {
		signedBody := append([]byte(clientUser.SessionKey), (*buf)[8:byteLength]...)

		generatedHash, deliveredHash := crypto.GenerateCRCHash(signedBody), binary.BigEndian.Uint32((*buf)[4:8])
		if generatedHash != deliveredHash {
			return nil, nil
		}
	}

	h.Lobby = clientLobby

	bd := (*buf)[20:byteLength]
	return &h, &bd
}

// TODO: assumption - user should be created on http requests not udp requests.

func UDPHandler(udpServer net.PacketConn, addr net.Addr, buf *[]byte, byteLength int) {
	header, body := parseHeader(buf, byteLength)
	// fmt.Printf("client user: %v\n", header)
	if header == nil {
		return
	}

	var response *[]byte
	doResponse := false
	var err error
	switch header.Command {
	case COMMAND_CONNECT:
		// do something
		// or does not happen?
	case COMMAND_JOIN:
		response, doResponse, err = onJoin(header, body)
	case COMMAND_PING:
		response, doResponse, err = onPing(header, body)
	default:
		// something is wrong. ignore packet
	}

	if err != nil {
		fmt.Printf("[%v] %v\n", time.Now().Format(time.ANSIC), err) // TODO: print timestamp to
	}

	if doResponse {
		responseBody := append([]byte{header.Command}, *response...)

		udpServer.WriteTo(responseBody, addr)
	}
}
