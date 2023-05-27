package lobby

import (
	maze "inagame/maze"
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
	Maze  maze.Maze
}

func NewLobby(id string, private bool) *Lobby {
	lobby := new(Lobby)
	lobby.Id = id
	lobby.Private = private
	lobby.Teams = []*Team{
		{Id: 0},
		{Id: 1},
	}

	return lobby
}

func (lobby *Lobby) InsertNewPlayer(id int64, name string) bool {
	inserted := false
	for _, team := range lobby.Teams {
		if len(team.Players) > 1 {
			continue
		}
		player := NewPlayer(id, name)
		if player == nil {
			return false
		}
		player.Team = team.Id
		player.IsReady = lobby.Private
		team.Players = append(team.Players, player)
		inserted = true
	}
	if !inserted {
		// todo: this shouldnt happen
	}
	return inserted
}
