package UDP

import "fmt"

func onJoin(header *Header, body *[]byte) (res *[]byte, reply bool, err error) {
	// TODO: parse bytes for username
	if header != nil && body != nil {
		fmt.Printf("%s joined lobby %s\n", header.User.Id, header.Lobby.Id)
		return &[]byte{0x00}, true, nil
	} else {
		return &[]byte{0x01}, true, nil
	}
}
