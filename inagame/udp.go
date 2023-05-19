package main

import (
	"fmt"
	"net"
	"time"
)

func UDPHandler(udpServer net.PacketConn, addr net.Addr, buf []byte, byteLength int) {
	time := time.Now().Format(time.ANSIC)
	responseStr := fmt.Sprintf("time received: %v. Your message: %v!", time, string(buf[0:byteLength]))
	fmt.Println(responseStr)
	// fmt.Println(len(string(buf)))

	udpServer.WriteTo([]byte(responseStr), addr)
}
