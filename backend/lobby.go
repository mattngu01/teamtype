// https://github.com/gorilla/websocket/tree/main/examples/chat

package backend

import "github.com/segmentio/ksuid"

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
		select {
		case client := <-l.register:
			l.clients[client] = true
		case client := <-l.unregister:
			delete(l.clients, client)
			// remember to close any client channels
		}
	}
}
