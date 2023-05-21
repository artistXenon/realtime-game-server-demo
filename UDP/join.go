package UDP

import (
	"encoding/json"
	"fmt"
)

type joinInfo struct {
	Id int
	GameId int
	Name string
}

func onJoin(body string) (res string, reply bool) {
	inputData := joinInfo{}
	err := json.Unmarshal([]byte(body), &inputData)
	if err != nil {
		fmt.Println(err)
		return "", false
	} else {
		fmt.Printf("%#v\n", inputData)
		return "ok", true // todo better message?
	}
}
