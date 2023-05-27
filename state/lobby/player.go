package lobby

// import "fmt"

type Player struct {
	Id   int64
	Name string

	LastPing     int64
	ReceiveDelay int16
	SendDelay    int16
	TimeOffset   int16
	Ping         int16

	// IsMe bool // is for client. not for server

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

var Players = make(map[int64]*Player)

func NewPlayer(id int64, name string) *Player {
	if Players[id] != nil {
		// should not happen
		return nil
	}
	p := new(Player)
	p.Id = id
	p.Name = name
	Players[id] = p
	return p
}

func DestroyPlayer(id int64) {
	// p := Players[id]
	// TODO: delete from lobby, delete from team
	delete(Players, id)
}
