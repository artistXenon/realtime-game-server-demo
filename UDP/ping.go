package UDP

import (
	// lobby "inagame/UDP/lobby"
	"encoding/binary"
	"errors"
	"time"
)

// client ping will include nothing but time

// receiveDelay: server time - sent client time
// sendDelay: sent server time - client time
// offset: server time - client time
// ping: average (server time - sent server time, client time - sent client time)

func onPing(header *Header, body *[]byte) (err error, res *[]byte, reply bool) {
	if header.User == nil {
		return errors.New("Failed to identify user on ping message"), nil, false
	}

	timeBytes := (*body)[0:8]
	sentTime := int64(binary.BigEndian.Uint64(timeBytes))

	lastPing := time.Now().UnixMilli()
	header.User.LastPing = lastPing
	receiveDelay := int16(lastPing - sentTime)
	header.User.ReceiveDelay = receiveDelay

	// generate ping res
	pingBytes := make([]byte, 10)
	binary.BigEndian.PutUint64(pingBytes, uint64(lastPing))
	binary.BigEndian.PutUint16(pingBytes[8:], uint16(receiveDelay))

	return nil, &pingBytes, true
}

func onPong(header *Header, body *[]byte) (err error, res *[]byte, reply bool) {
	if header.User == nil {
		return errors.New("Failed to identify user on ping message"), nil, false
	}

	_, offset, sendDelay := binary.BigEndian.Uint16(*body), binary.BigEndian.Uint16((*body)[2:]), binary.BigEndian.Uint16((*body)[4:])

	header.User.Ping = int16(time.Now().UnixMilli() - header.User.LastPing)
	header.User.SendDelay = int16(sendDelay)
	header.User.TimeOffset = int16(offset)

	return nil, nil, false
}
