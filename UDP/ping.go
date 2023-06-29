package UDP

import (
	// lobby "inagame/UDP/lobby"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// client ping will include nothing but time

// receiveDelay: server time - sent client time
// sendDelay: sent server time - client time
// offset: server time - client time
// ping: average (server time - sent server time, client time - sent client time)

func onPing(header *Header, body *[]byte) (res *[]byte, reply bool, err error) {
	if header.User == nil {
		return nil, false, errors.New("failed to identify user on ping message")
	}

	timeBytes := (*body)[0:8]
	clientTime, sendDelay := int64(binary.BigEndian.Uint64(timeBytes)), int16(binary.BigEndian.Uint16((*body)[8:]))

	now := time.Now().UnixMilli()
	receiveDelay := int16(now - clientTime)

	var ping, offset int16 = 0, 0
	if sendDelay != -32768 {
		ping = int16(now - header.User.LastPing)
		offset = (int16(sendDelay) + receiveDelay) / 2
		header.User.AppendPing(ping, offset)
	}

	header.User.LastPing = now
	header.User.ReceiveDelay = receiveDelay

	// generate ping res
	pingBytes := make([]byte, 10)
	binary.BigEndian.PutUint64(pingBytes, uint64(now))
	binary.BigEndian.PutUint16(pingBytes[8:], uint16(receiveDelay))
	header.Command = COMMAND_PONG

	fmt.Printf("ping: %d offset: %d\n", header.User.AvgPing(), header.User.AvgOffset())
	return &pingBytes, true, nil
}
