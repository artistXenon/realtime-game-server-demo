package HTTP

import (
	"fmt"
	"inagame/db"
	"inagame/state"
	"inagame/state/lobby"
	"net/http"
	"strconv"
	"sync"
)

// type httpResponse struct {
// 	server string
// 	lobby  string
// }

// TODO: send status code on Write
func HTTPCreateHandler(w http.ResponseWriter, req *http.Request) {
	// http://url?id=fwp3js

	ids, iok := req.URL.Query()["id"]
	id := ids[0]
	prvs, pok := req.URL.Query()["private"]
	prv := prvs[0]

	if !iok || len(id) < 1 || !pok || len(prv) < 1 {
		w.WriteHeader(400)
		w.Write([]byte("wrong query"))
		return
	}

	mtx := new(sync.RWMutex)
	mtx.Lock()
	defer mtx.Unlock()

	if len(lobby.Lobbys) >= state.GameCapacity {
		w.WriteHeader(400)
		w.Write([]byte("no capacity"))
		return
	}
	b, err := strconv.ParseBool(prv)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("wrong format"))
		return
	}
	lb := lobby.NewLobby(id, b)
	lobby.Lobbys[id] = lb

	fmt.Printf("lobby was created with id %s\n", id)

	w.Write([]byte("lobby created"))
	// todo: send server address
}

func HTTPJoinHandler(w http.ResponseWriter, req *http.Request) {
	ids, iok := req.URL.Query()["uid"]
	uid := ids[0]
	lids, lok := req.URL.Query()["lid"]
	lid := lids[0]

	if !iok || len(uid) < 1 {
		w.Write([]byte("no user id"))
		return
	}
	if !lok || len(lid) < 1 {
		w.Write([]byte("no lobby id"))
		return
	}

	mtx := new(sync.RWMutex)
	mtx.Lock()
	defer mtx.Unlock()

	clientLobby := lobby.Lobbys[lid]
	clientUser := lobby.Players[uid]

	if clientLobby == nil { // shouldn't happen
		w.WriteHeader(400)
		w.Write([]byte("lobby does not exist"))
		return
	}

	if clientUser == nil {
		sessionKey := db.GetPlayer(uid, lid)
		if sessionKey == nil {
			w.WriteHeader(400)
			w.Write([]byte("user does not exist"))
			return
		}
		lobby.CreatePlayer(uid, *sessionKey, clientLobby)
		// previous player joined new lobby. re assign player w/ refreshed session key
	} else if clientUser.Lobby.Id != clientLobby.Id {
		sessionKey := db.GetPlayer(uid, lid)
		if sessionKey == nil {
			w.WriteHeader(400)
			w.Write([]byte("user does not exist"))
			return
		}
		clientUser.SessionKey = *sessionKey
		clientLobby.AssignPlayer(clientUser)
	} else {
		sessionKey := db.GetPlayer(uid, lid)
		clientUser.SessionKey = *sessionKey
	}
	w.Write([]byte("user joined"))
	// create user by db fetch.
	// clean users who were previously present in different lobby.
	// clear user if after http join, udp join does not occur within certain amount of time.
}
