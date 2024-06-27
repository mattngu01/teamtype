package main

import (
	"log"
	"net/http"

	"backend/backend"

	"github.com/segmentio/ksuid"
)

var addr string = ":8080"

func main() {
	log.Println("Starting up server")
	lm := &backend.LobbyManager{make(map[ksuid.KSUID]*backend.Lobby)}

	http.HandleFunc("/ws", lm.ServeWs)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
