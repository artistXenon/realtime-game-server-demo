package TCP

import (
	"encoding/binary"
	"inagame/state/lobby"
)

func lobbyHandler(buf *[]byte, player *lobby.Player) (res *[]byte, disconnect bool, err error) {
	resBytes := &[]byte{(byte)(player.Lobby.State)}

	for _, team := range player.Lobby.Teams {
		missing := 2
		for _, p := range team.Players {
			playerByte := make([]byte, 30)
			playerByte[0] = 0b1000_0000
			// if p.IsReady {
			// 	playerByte[0] |= 0b0100_0000
			// }
			if p.IsLeader {
				playerByte[0] |= 0b0010_0000
			}
			if p.Id == player.Id {
				playerByte[0] |= 0b0001_0000
			}

			playerByte[1] = byte(p.Team)

			copy(playerByte[2:], []byte(p.Name))
			copy(playerByte[18:], []byte(p.InternalId))
			binary.BigEndian.PutUint16(playerByte[23:], 0)
			*resBytes = append(*resBytes, playerByte...)
			missing -= 1
		}
		for i := 0; i < missing; i++ {
			playerByte := make([]byte, 30)
			playerByte[0] = 0b0000_0000
			*resBytes = append(*resBytes, playerByte...)
		}
	}

	return resBytes, false, nil
}

/*
state 1

player1 state 1 []
exist
ready
leader
isme

player1 team 1
player1 name 16
player1 id 5
player1 character 2
*/

func lobbyLeaveHandler(buf *[]byte, player *lobby.Player) (res *[]byte, disconnect bool, err error) {
	resBytes := &[]byte{0x00}
	player.Lobby.RemovePlayer(player)
	// MAYBE TODO: update db so that player has no lobby
	return resBytes, true, nil
}

// status code of leave.
/*
0: requested leave
1: lobby kicked
2: server's decision (bad internet)
*/
