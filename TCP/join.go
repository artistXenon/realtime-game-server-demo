package TCP

import (
	"encoding/binary"
	"inagame/crypto"
	"inagame/state/lobby"
	"math/big"
)

func joinHandler(buf *[]byte, length int) (player *lobby.Player, errorCode int8) {
	idInt := new(big.Int)
	idInt.SetBytes((*buf)[5:15])
	idString := idInt.String()
	lobbyString := string((*buf)[15:20])
	// TODO: name will be included

	clientUser := lobby.Players[idString]
	clientLobby := lobby.Lobbys[lobbyString]

	// new player on this server. create and assign
	if clientUser == nil || clientLobby == nil || clientUser.Lobby.Id != clientLobby.Id {
		return nil, 1
		// previous player joined new lobby. re assign player w/ refreshed session key
	}

	signedBody := append([]byte(clientUser.SessionKey), (*buf)[4:length]...)

	generatedHash := crypto.GenerateCRCHash(signedBody)
	if generatedHash != binary.BigEndian.Uint32((*buf)[0:4]) {
		return nil, 2
	}

	return clientUser, 0
}
