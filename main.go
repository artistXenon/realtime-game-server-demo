package main

import (
	"fmt"
	"log"
	"net"
    "net/http"

	globalContext "inagame/global-context"
	ina_http "inagame/HTTP"
	ina_udp "inagame/UDP"

	// mazer "inagame/maze"
)


func main() {
	ch := make(chan int)

	// maze := mazer.NewMaze()
	// maze.SetWidth(5).SetHeight(5).Init()

	// go maze.Generate(ch)
	fmt.Println(globalContext.GameCapacity)
	go startHTTP()
	go startUDP()

	x := <- ch
	fmt.Println(x)
	// fmt.Printf(maze.Serialize())
}	


func startHTTP() {
    http.HandleFunc("/hello", ina_http.HTTPHandler)
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
			go ina_udp.UDPHandler(udpServer, addr, buf, resLen)
		}
	}
}