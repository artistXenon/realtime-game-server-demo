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

func (lobby *Lobby) AssignPlayer(player *Player) bool {
	player.Lobby.RemovePlayer(player)

	for _, team := range lobby.Teams {
		if len(team.Players) > 1 {
			continue
		}
		player.Lobby = lobby
		player.Team = team.Id
		player.IsReady = lobby.Private
		team.Players = append(team.Players, player)
		return true
	}

	// TODO: one can reach here when insert failed. this shouldnt happen

	return false
}

func (lobby *Lobby) RemovePlayer(player *Player) bool {
	if player.Lobby == nil || player.Lobby.Id != lobby.Id {
		return false
	}
	for _, team := range lobby.Teams {
		for teamIndex, teamPlayer := range team.Players {
			if teamPlayer.Id == player.Id {
				team.Players = append(team.Players[:teamIndex], team.Players[teamIndex+1:]...)
				return true
			}
		}
	}
	return false
}

// TODO: rewrite inset new player w/ arg: Player
// needs to check if player exist.
// if exists, remove player from previous lobby
// make sure this is clean.
// and then insert this player to a `this` lobby
