package UDP

import (
	"encoding/json"
	"fmt"
	"inagame/state"
	"sync"
)

type joinInfo struct { // JSON
	GameId string
	Name   string
}

func onJoin(msg *Message) (res string, reply bool) {
	mtx := new(sync.RWMutex) //TODO: fix sync. this just doesn't feel right to make new lock every time a request comes.
	inputData := joinInfo{}
	err := json.Unmarshal([]byte(msg.Body), &inputData)
	if err != nil {
		fmt.Println(err)
		return "!", false
	}

	fmt.Println(inputData)

	// var game *lobby.Lobby
	mtx.RLock()
	game := state.Games[inputData.GameId]
	mtx.RUnlock()

	if game == nil {
		return "!wrong server", true
	}

	mtx.Lock()
	defer mtx.Unlock()

	game.InsertNewPlayer(msg.UserId, inputData.Name)

	fmt.Printf("%#v\n", inputData)
	return "ok", true // todo better message?
}
