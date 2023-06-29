package lobby

import (
	"net"
)

// import "fmt"

const ping_buffer = 20

type Player struct {
	Id         string
	Name       string
	SessionKey string
	TCP        *net.Conn

	LastPing     int64
	ReceiveDelay int16
	SendDelay    int16
	TimeOffset   []int16
	Ping         []int16

	Lobby    *Lobby
	JoinTime int64

	// mutable on waiting lobby
	Team     int8
	IsLeader bool
	IsReady  bool

	// mutable on gaming lobby
	PositionX int16
	PositionY int16

	// not implemented
	Cosmetics Cosmetics
}

var Players = make(map[string]*Player)

func CreatePlayer(id string, sessionKey string, lobby *Lobby) (player *Player) {
	p := &Player{
		Id:         id,
		SessionKey: sessionKey,
	}
	lobby.AssignPlayer(p)
	Players[id] = p
	return p
}

func DestroyPlayer(id string) {
	player := Players[id]
	player.Lobby.RemovePlayer(player)
	delete(Players, id)
}

func (player *Player) AppendPing(ping int16, offset int16) {
	player.Ping = append(player.Ping, ping)
	player.TimeOffset = append(player.TimeOffset, offset)

	pingOverSize := len(player.Ping) - ping_buffer
	if pingOverSize > 0 {
		player.Ping = player.Ping[pingOverSize:]
	}
	offsetOverSize := len(player.TimeOffset) - ping_buffer
	if offsetOverSize > 0 {
		player.TimeOffset = player.TimeOffset[offsetOverSize:]
	}
}

func (player *Player) AvgPing() (ping int16) {
	pingsLength := int16(len(player.Ping))
	if pingsLength == 0 {
		return 0
	}
	sumPing := int16(0)
	for _, v := range player.Ping {
		sumPing += v
	}
	return sumPing / pingsLength
}

func (player *Player) AvgOffset() (offset int16) {
	offsetsLength := int16(len(player.TimeOffset))
	if offsetsLength == 0 {
		return 0
	}
	sumOffset := int16(0)
	for _, v := range player.TimeOffset {
		sumOffset += v
	}
	return sumOffset / offsetsLength
}
