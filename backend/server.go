package main

import (
	"log"
	"net/http"

	"github.com/segmentio/ksuid"
)

var addr string = ":8080"

func main() {
	log.Println("Starting up server")
	lm := &LobbyManager{make(map[ksuid.KSUID]*Lobby)}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, lm)
	})
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
