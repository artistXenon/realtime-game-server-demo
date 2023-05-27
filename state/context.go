package state

import lobby "inagame/state/lobby"

var GameCapacity int = 10

var Games []*lobby.Lobby = make([]*lobby.Lobby, GameCapacity)
