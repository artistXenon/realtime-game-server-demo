package state

import "inagame/state/lobby"

var GameCapacity int = 10

// todo: make Games a map[string]Lobby
var Games = make(map[string]*lobby.Lobby)
