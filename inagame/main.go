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
		buf := make([]byte, 1024)
		_, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		go UDPHandler(udpServer, addr, buf)
	}
}