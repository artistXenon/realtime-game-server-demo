package main

import (
	"fmt"
	"log"
	"net"
    "net/http"
)


func main() {
	ch := make(chan int)
	go startHTTP()
	go startUDP()

	x := <- ch
	fmt.Println(x)
}	


func startHTTP() {
    http.HandleFunc("/hello", HTTPHandler)

    http.ListenAndServe(":5000", nil)
}

func startUDP() {
	udpServer, err := net.ListenPacket("udp", ":5001")
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()

	for {
		// todo:
		// make byte array pool for each goroutine to make use of and return for reuse.
		// stream of long packet may be dealt with byte arrays for each streams
		buf := make([]byte, 1024)
		resLen, addr, err := udpServer.ReadFrom(buf)
		if resLen > 0 {
			if err != nil {
				continue
			}
			go UDPHandler(udpServer, addr, buf, resLen)
		}
	}
}