package lobby

type Team struct {
	Players []*Player
	// Items   []int8  // TODO: create item type
	Gauge float32 // team associated gauge for ult.
}
