package main

import (
	"fmt"
	"net"
	"time"
)

func UDPHandler(udpServer net.PacketConn, addr net.Addr, buf []byte) {
	time := time.Now().Format(time.ANSIC)
	responseStr := fmt.Sprintf("time received: %v. Your message: %v!", time, string(buf))
	fmt.Println(responseStr)

	udpServer.WriteTo([]byte(responseStr), addr)
}
