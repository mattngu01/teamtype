package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// acts as middleman between frontend websocket and a Lobby
type Client struct {
	socket *websocket.Conn
	input  chan string
	lobby  *Lobby
}

/*
TODO:
- going through workflow of creating/joining a lobby
- frontend will establish websocket, backend create client obj & lobby respectively
- there needs to be a server listening for websocket requests? and to create the lobbies
- make sure backend acknowledges leavers
- after we can move onto game state
*/
