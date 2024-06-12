package backend

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

func serveWs(w http.ResponseWriter, r *http.Request, l *Lobby) {
	conn, err := upgrader.Upgrade(w, r, nil)

	client := &Client{socket: conn, input: make(chan string), lobby: l}
	client.lobby.register <- client
	/*
		go run readRoutine
		go run writeRoutine
	*/
}

/*
TODO:
- going through workflow of creating/joining a lobby
- frontend will establish websocket, backend create client obj & lobby respectively
- there needs to be a server listening for websocket requests? and to create the lobbies
- make sure backend acknowledges leavers
- after we can move onto game state
*/
