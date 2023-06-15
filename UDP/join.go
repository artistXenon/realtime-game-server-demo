package UDP

type joinInfo struct { // JSON
	LobbyId string
	Name    string
}

func onJoin(header *Header, body *[]byte) (err error, res *[]byte, reply bool) {
	if header != nil && body != nil {
		return nil, &[]byte{0x00}, true
	} else {
		return nil, &[]byte{0x01}, true
	}
}
