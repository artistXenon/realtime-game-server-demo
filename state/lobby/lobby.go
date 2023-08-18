package lobby

import (
	maze "inagame/maze"
	"time"
)

const ( // enum Lobby State
	Waiting = iota + 1
	Gaming
)

type Lobby struct {
	Id      string
	Private bool
	State   int8 // enum Lobby State

	Teams []*Team
	Maze  *maze.Maze
}

var Lobbys = make(map[string]*Lobby)

func NewLobby(id string, private bool) *Lobby {
	// TODO: check if id already exists in match list
	lobby := new(Lobby)
	lobby.Id = id
	lobby.Private = private
	team1 := &Team{Players: []*Player{}}
	team2 := &Team{Players: []*Player{}}
	lobby.Teams = []*Team{
		team1, team2,
	}

	return lobby
}

func (lobby *Lobby) DestroyLobby(broadcast bool) {
	if broadcast {
		// disconnect connected players
	}
	delete(Lobbys, lobby.Id)
	// remove from db

}

func (lobby *Lobby) GetAllPlayers() []*Player {
	playersCache := []*Player{}
	for _, team := range lobby.Teams {
		for _, tPlayer := range team.Players {
			playersCache = append(playersCache, tPlayer)
		}
	}
	return playersCache
}

func (lobby *Lobby) AssignPlayer(player *Player) bool {
	player.Lobby.RemovePlayer(player)

	// TODO: check if any player is out dated after join

	created := false
	leader := false
	for id, team := range lobby.Teams {
		teamSize := len(team.Players)
		var i int
		for i = 0; i < teamSize; i++ {
			if team.Players[i].IsLeader {
				leader = true
				break
			}
		}
		if i == 2 {
			continue
		}

		if created {
			continue
		}
		player.JoinTime = time.Now().UnixMilli()
		player.Lobby = lobby
		player.Team = int8(id)
		// player.IsReady = !lobby.Private
		team.Players = append(team.Players, player)
		created = true
	}
	if created && !leader {
		player.IsLeader = true
	}
	// one can reach here when insert failed. this shouldn't happen
	return created
}

func (lobby *Lobby) RemovePlayer(player *Player) bool {
	if player.Lobby == nil || player.Lobby.Id != lobby.Id {
		return false
	}

	removed := false
	leaderTransferred := !player.IsLeader
	playersCache := lobby.GetAllPlayers()
	for _, cPlayer := range playersCache {
		if cPlayer.Id == player.Id {
			team := lobby.Teams[player.Team]
			for tid, tPlayer := range team.Players {
				if cPlayer.Id == tPlayer.Id {
					team.Players = append(team.Players[:tid], team.Players[tid+1:]...)
					removed = true
				}
			}
		} else if !leaderTransferred {
			cPlayer.IsLeader = true
			leaderTransferred = true
		}
	}

	if !leaderTransferred {
		lobby.DestroyLobby(false)
	} else {
		// TODO
		// if lobby state is waiting
		// for pid, tPlayer := range playersCache {
		// 	// (*tPlayer.TCP).Write()
		// 	// all players receive lobby update
		// }
	}

	return removed
}
