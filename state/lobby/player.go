package lobby

// import "fmt"

type Player struct {
	Id         string
	Name       string
	SessionKey string

	LastPing     int64
	ReceiveDelay int16
	SendDelay    int16
	TimeOffset   int16
	Ping         int16

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
