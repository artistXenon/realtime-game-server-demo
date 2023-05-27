package HTTP

import (
	// "fmt"
	"inagame/state"
	"inagame/state/lobby"
	"net/http"
	"strconv"
	"sync"
	// "time"
)

func HTTPHandler(w http.ResponseWriter, req *http.Request) {
	// http://url?id=fwp3js

	ids, ok := req.URL.Query()["id"]
	prv, ok := req.URL.Query()["private"]

	if !ok || len(ids[0]) < 1 {
		w.Write([]byte("no id")) // todo: send status code
		return
	}

	idEmpty := -1
	mtx := new(sync.RWMutex)
	mtx.Lock()
	defer mtx.Unlock()
l1:
	for index, element := range state.Games {
		if element == nil {
			idEmpty = index
			break l1
		}
	}

	if idEmpty == -1 {
		w.Write([]byte("no capacity")) // todo: send status code
	} else {
		b, err := strconv.ParseBool(prv[0])
		if err != nil {
			w.Write([]byte("wrong format"))
			return
		}
		state.Games[idEmpty] = lobby.NewLobby(ids[0], b)
		w.Write([]byte("match created"))
		// todo: send status code
		// todo: send server address
		// strconv.FormatInt(time.Now().UnixMilli(), 10)
	}
}
