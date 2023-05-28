package UDP

import (
	// lobby "inagame/UDP/lobby"
	"encoding/json"
	"fmt"
	"inagame/state/lobby"
	"strconv"
	"time"
)

type times struct { // this sucks more than what I expected
	PlayerId  string `json:",omitempty"`
	Delay     int64
	LocalTime int64
	Offset    int64
	Ping      int64
}

type pong struct {
	Ping      int16
	Offset    int16
	SendDelay int16
}

// client ping will include nothing but time

// receiveDelay: server time - sent client time
// sendDelay: sent server time - client time
// offset: server time - client time
// ping: average (server time - sent server time, client time - sent client time)

func onPing(msg *Message) (res string, reply bool) {
	clientTime, error := strconv.ParseInt(msg.Body, 0, 64) //TODO: json parse this thing.
	if error != nil {
		// wrong stuff
		return "!not client", false
	}
	player := lobby.Players[msg.UserId]

	if player == nil {
		return "!not client", false
	}
	player.LastPing = time.Now().UnixMilli()

	player.ReceiveDelay = int16(player.LastPing - clientTime) // ping + offset

	// todo: record this info for client
	// player := nil //*lobby.Player     <-- definition required before unix call
	return fmt.Sprintf("ping!%d!%d!", player.LastPing, player.ReceiveDelay), true

}

func onPong(msg *Message) {
	player := lobby.Players[msg.UserId]
	if player == nil {
		return
	}
	p := pong{}
	json.Unmarshal([]byte(msg.Body), &p)

	player.Ping = int16(time.Now().UnixMilli() - player.LastPing)
	player.SendDelay = p.SendDelay
	player.TimeOffset = p.Offset

	fmt.Printf("client calculated: %d %d\n", p.Offset, p.Ping)

	fmt.Printf("server calculated: %d %d\n", p.Offset, player.Ping)

	// todo: do something with tTimes
}
