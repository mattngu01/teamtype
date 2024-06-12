package main

// https://github.com/gorilla/websocket/tree/main/examples/chat

import (
	"log"
	"net/http"

	"github.com/segmentio/ksuid"
)

type LobbyManager struct {
	lobbies map[ksuid.KSUID]*Lobby
}

func serveWs(w http.ResponseWriter, r *http.Request, lm *LobbyManager) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	// is it okay for each client to have their own lobby on initial connection..?
	lobby := newLobby()
	go lobby.run()
	client := &Client{socket: conn, input: make(chan string), lobby: lobby}
	lm.lobbies[lobby.id] = lobby
	log.Println("Created client & lobby")
	log.Println("Sending register client msg to lobby", client.lobby.id)
	client.lobby.register <- client
}

type Lobby struct {
	race       *Race
	clients    map[*Client]bool
	id         ksuid.KSUID
	register   chan *Client
	unregister chan *Client
}

func newLobby() *Lobby {
	return &Lobby{
		race:       newRace(getQuote()),
		clients:    make(map[*Client]bool, 8),
		id:         ksuid.New(),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func getQuote() string {
	return "We don't understand what really causes events to happen. History is the fiction we invent to persuade ourselves that events are knowable and that life has order and direction. That's why events are always reinterpreted when values change. We need new versions of history to allow for our current prejudices."
}

func (l *Lobby) run() {
	for {
		if len(l.clients) == 0 {
			return
		}

		select {
		case client := <-l.register:
			log.Println("Registered new client")
			l.clients[client] = true
		case client := <-l.unregister:
			log.Println("Deregistered client")
			delete(l.clients, client)
			// remember to close websocket, as well as close goroutine channel to signal cleaning up goroutine
			close(client.input)
		}
	}
}
