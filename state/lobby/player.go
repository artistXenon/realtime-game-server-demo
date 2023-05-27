package lobby

type Player struct {
	Id   string
	Name string
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

var Players = make(map[string]*Player)

func NewPlayer(id string, name string) *Player {
	p := new(Player)
	p.Id = id
	p.Name = name
	Players[id] = p
	return p
}

func DestroyPlayer(id string) {
	// p := Players[id]
	// TODO: delete from lobby, delete from team
	delete(Players, id)
}
