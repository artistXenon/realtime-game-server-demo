package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"inagame/HTTP"
	"inagame/TCP"
	"inagame/UDP"
	"inagame/state"
	// mazer "inagame/maze"
)

func main() {
	ch := make(chan int)

	// maze := mazer.NewMaze()
	// maze.SetWidth(5).SetHeight(5).Init()

	// go maze.Generate(ch)
	fmt.Println(state.GameCapacity)
	go startHTTP()
	go startUDP()
	go startTCP()

	x := <-ch
	fmt.Println(x)
	// fmt.Printf(maze.Serialize())
}

func startHTTP() {
	http.HandleFunc("/create", HTTP.HTTPCreateHandler)
	http.HandleFunc("/join", HTTP.HTTPJoinHandler)
	http.ListenAndServe(":5000", nil) // TODO: remove localhost
}

func startUDP() {
	udpServer, err := net.ListenPacket("udp", ":5001") // TODO: remove localhost
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
			go UDP.UDPHandler(udpServer, addr, &buf, resLen)
		}
	}
}

func startTCP() {
	tcpServer, err := net.Listen("tcp", ":5003")
	if err != nil {
		log.Fatal(err)
	}
	defer tcpServer.Close()

	for {
		conn, err := tcpServer.Accept()
		if err != nil {
			continue
		}
		go TCP.TCPHandler(conn)
	}
}
