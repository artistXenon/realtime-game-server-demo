package main

import (
	"fmt"
	"net/http"
)

func HTTPHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello World"))
	fmt.Println("foooo!")
}