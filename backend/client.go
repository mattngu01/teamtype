package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// acts as middleman between frontend websocket and a Lobby
type Client struct {
	conn *websocket.Conn
	lobbyEvents  chan Event //channel to pass events from client to lobby, where it maintains game state
	lobby  *Lobby
	clientEvents chan Event //channel to pass events from writer & reader goroutines, skipping Lobby when not needed
	username string //usernames unique per lobby
}

func newClient(conn *websocket.Conn, lobby *Lobby) *Client {
	return &Client{conn: conn, lobbyEvents: make(chan Event), lobby: lobby, clientEvents: make(chan Event), username: petname.Generate(3, "-")}
}

// at most one reader on a connection by executing all reads on this goroutine
func (c *Client) readRoutine() {
	log.Println("Starting read routine")
	defer func() {
		log.Println("Read routine, closing conn")
		c.lobby.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		event := &Event{}
		if err = json.Unmarshal(message, event); err != nil {
			log.Printf("Unable to unmarshal msg: %v", err)
		}

		c.parseEvent(event)
	}
}

func (c *Client) parseEvent(event *Event) {
	log.Printf("Parsing event: %v", event)
	switch event.Type {
	case LobbyInfo:
		c.clientEvents <- Event{Type: LobbyInfo, Data: map[string]interface{}{"lobbyId": c.lobby.id.String(), "username": c.username, "players": c.lobby.getPlayerUsers()}}
	}
}

func (c *Client) writeRoutine() {
	log.Println("Starting write routine")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("Write routine, closing conn")
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Failure to write message %v", err)
				return
			}
		case event := <- c.clientEvents:
			if (event.Type == LobbyInfo) {
				w, err := c.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					log.Printf("Error obtaining writer: %v", err)
					return
				}

				msg, err := json.Marshal(event)
				if err != nil {
					log.Printf("Unable to marshal message: %v", err)
					break
				}
				w.Write(msg)

				if err = w.Close(); err != nil {
					log.Printf("Unable to close writer: %v", err)
					return
				}
			}
		}
	}
}

/*
TODO:
- going through workflow of creating/joining a lobby
- frontend will establish websocket, backend create client obj & lobby respectively
- there needs to be a server listening for websocket requests? and to create the lobbies
- make sure backend acknowledges leavers
- after we can move onto game state
*/
