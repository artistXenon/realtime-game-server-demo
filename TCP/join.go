package TCP

import (
	"encoding/binary"
	"errors"
	"inagame/crypto"
	"inagame/state"
	"inagame/state/lobby"
	"math/big"
)

func joinHandler(buf *[]byte, length int) (player *lobby.Player, err error) {
	idInt := new(big.Int)
	idInt.SetBytes((*buf)[5:15])
	idString := idInt.String()
	lobbyString := string((*buf)[15:20])
	// TODO: name will be included

	clientUser := lobby.Players[idString]
	clientLobby := state.Games[lobbyString]

	// new player on this server. create and assign
	if clientUser == nil || clientLobby == nil || clientUser.Lobby.Id != clientLobby.Id {
		return nil, errors.New("matching user/lobby not found")
		// previous player joined new lobby. re assign player w/ refreshed session key
	}

	signedBody := append([]byte(clientUser.SessionKey), (*buf)[4:length]...)

	generatedHash := crypto.GenerateCRCHash(signedBody)
	if generatedHash != binary.BigEndian.Uint32((*buf)[0:4]) {
		return nil, errors.New("invalid hash on join request")
	}

	return clientUser, nil
}
