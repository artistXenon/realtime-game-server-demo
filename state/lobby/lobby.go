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

func NewLobby(id string, private bool) *Lobby {
	// TODO: check if id already exists in match list
	lobby := new(Lobby)
	lobby.Id = id
	lobby.Private = private
	lobby.Teams = []*Team{
		new(Team),
		new(Team),
	}

	return lobby
}

func (lobby *Lobby) AssignPlayer(player *Player) bool {
	player.Lobby.RemovePlayer(player)

	// TODO: check if any player is out dated after join

	for id, team := range lobby.Teams {
		if len(team.Players) > 1 {
			continue
		}
		player.JoinTime = time.Now().UnixMilli()
		player.Lobby = lobby
		player.Team = int8(id)
		player.IsReady = lobby.Private
		team.Players = append(team.Players, player)
		return true
	}

	// one can reach here when insert failed. this shouldn't happen

	return false
}

func (lobby *Lobby) RemovePlayer(player *Player) bool {
	if player.Lobby == nil || player.Lobby.Id != lobby.Id {
		return false
	}
	team := lobby.Teams[player.Team]

	for teamIndex, teamPlayer := range team.Players {
		if teamPlayer.Id == player.Id {
			team.Players = append(team.Players[:teamIndex], team.Players[teamIndex+1:]...)
			return true
		}
	}

	return false
}
