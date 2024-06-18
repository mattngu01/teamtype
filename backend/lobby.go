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
	client := newClient(conn)
	lobby := newLobby(client)
	client.lobby = lobby
	go lobby.run()
	lm.lobbies[lobby.id] = lobby
	log.Println("Created client & lobby")

	go client.readRoutine()
	go client.writeRoutine()
}

type Lobby struct {
	race       *Race //represents game state
	clients    map[*Client]bool
	id         ksuid.KSUID
	register   chan *Client
	unregister chan *Client
}

// requires first client since run() goroutine exits if there aren't any clients
func newLobby(initialClient *Client) *Lobby {
	clients := make(map[*Client]bool, 8)
	clients[initialClient] = true
	return &Lobby{
		race:       newRace(getQuote()),
		clients:    clients,
		id:         ksuid.New(),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (l *Lobby) getPlayerUsers() []string {
	var usernames []string
	for clientPtr := range l.clients {
		usernames = append(usernames, clientPtr.username)
	}
	return usernames
}

func getQuote() string {
	return "We don't understand what really causes events to happen. History is the fiction we invent to persuade ourselves that events are knowable and that life has order and direction. That's why events are always reinterpreted when values change. We need new versions of history to allow for our current prejudices."
}

func (l *Lobby) run() {
	for {
		if len(l.clients) == 0 {
			log.Println("Killing lobby routine, 0 clients")
			return
		}

		select {
		case client := <-l.register:
			log.Println("Registered new client", client.username)
			l.clients[client] = true
		case client := <-l.unregister:
			log.Println("Deregistered client", client.username)
			delete(l.clients, client)
			// remember to close websocket, as well as close goroutine channel to signal cleaning up goroutine
			close(client.lobbyEvents)
		}
	}
}
