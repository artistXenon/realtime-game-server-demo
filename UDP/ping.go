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

	prevTick := header.User.UDPT
	curTick := binary.BigEndian.Uint32(header.Count)
	lostTicks := int16(curTick - prevTick - 1)
	header.User.AppendLoss(lostTicks)

	timeBytes := (*body)[0:8]
	clientTime, sendDelay := int64(binary.BigEndian.Uint64(timeBytes)), int16(binary.BigEndian.Uint16((*body)[8:]))

	now := time.Now().UnixMilli()
	receiveDelay := int16(now - clientTime)

	var ping, offset int16 = 0, 0
	if sendDelay != -32768 {
		ping = int16(now - header.User.LastPing)
		offset = (sendDelay + receiveDelay) / 2
		header.User.AppendPing(ping, offset)
	} else {

	}

	header.User.LastPing = now
	header.User.ReceiveDelay = receiveDelay
	header.User.UDPT = curTick

	// generate ping res
	pingBytes := make([]byte, 10)
	binary.BigEndian.PutUint64(pingBytes, uint64(now))
	binary.BigEndian.PutUint16(pingBytes[8:], uint16(receiveDelay))
	header.Command = COMMAND_PONG

	fmt.Printf("ping: %d offset: %d loss: %f\n", header.User.AvgPing(), header.User.AvgOffset(), header.User.LossRate())
	return &pingBytes, true, nil
}
