package HTTP

import (
	"fmt"
	"inagame/state"
	"inagame/state/lobby"
	"net/http"
	"strconv"
	"sync"
	// "time"
)

func HTTPHandler(w http.ResponseWriter, req *http.Request) {
	// http://url?id=fwp3js

	ids, iok := req.URL.Query()["id"]
	id := ids[0]
	prvs, pok := req.URL.Query()["private"]
	prv := prvs[0]

	if !iok || len(id) < 1 {
		w.Write([]byte("no id")) // todo: send status code
		return
	}
	if !pok || len(prv) < 1 {
		w.Write([]byte("no accessor")) // todo: send status code
		return
	}

	mtx := new(sync.RWMutex)
	mtx.Lock()
	defer mtx.Unlock()

	if len(state.Games) >= state.GameCapacity {
		w.Write([]byte("no capacity")) // todo: send status code
	} else {
		b, err := strconv.ParseBool(prv)
		if err != nil {
			w.Write([]byte("wrong format"))
			return
		}
		// todo: check if id already exists in match list
		lb := lobby.NewLobby(id, b)
		state.Games[id] = lb

		fmt.Printf("%#v\n", state.Games)
		w.Write([]byte("match created"))
		// todo: send status code
		// todo: send server address
		// strconv.FormatInt(time.Now().UnixMilli(), 10)
	}
}
