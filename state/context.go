package state

import "inagame/state/lobby"

var GameCapacity int = 10

var Games = make(map[string]*lobby.Lobby)
