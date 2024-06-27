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

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		backend.ServeWs(w, r, lm)
	})
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
