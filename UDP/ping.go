package UDP

import (
	// lobby "inagame/UDP/lobby"
	"encoding/json"
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

// receiveDelay: server time - sent client time
// sendDelay: sent server time - client time
// offset: server time - client time
// ping: average (server time - sent server time, client time - sent client time)

func onPing(body string) (res string, reply bool) {
	clientTime, error := strconv.ParseInt(body, 0, 64) //TODO: json parse this thing.
	if error != nil {
		// wrong stuff
		return "not client", false
	}
	tTimes := times{}
	tTimes.LocalTime = time.Now().UnixMilli()

	tTimes.Delay = tTimes.LocalTime - clientTime // ping + offset

	jsonTimes, _ := json.Marshal(tTimes) // err should never happens

	// todo: record this info for client
	// player := nil //*lobby.Player     <-- definition required before unix call
	return string(jsonTimes), true

}

func onPong(body string) {
	tTimes := times{}
	json.Unmarshal([]byte(body), &tTimes)

	// todo: do something with tTimes
}
