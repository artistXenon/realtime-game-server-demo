package UDP

import (
	"encoding/json"
	"fmt"
	"inagame/state"
	"sync"
)

type joinInfo struct { // JSON
	Id     string
	GameId string
	Name   string
}

func onJoin(body string) (res string, reply bool) {
	mtx := new(sync.RWMutex)
	inputData := joinInfo{}
	err := json.Unmarshal([]byte(body), &inputData) // aware: id in multiple game
	if err != nil {
		fmt.Println(err)
		return "!", false
	}

	gameIndex := -1
	mtx.RLock()
l1:
	for index, element := range state.Games {
		if inputData.GameId == element.Id {
			gameIndex = index
			break l1
		}
	}
	mtx.Unlock()

	if gameIndex == -1 {
		return "!wrong server", true
	}

	mtx.Lock()
	game := state.Games[gameIndex]

	game.InsertNewPlayer(inputData.Id, inputData.Name)

	fmt.Printf("%#v\n", inputData)
	return "ok", true // todo better message?
}
