package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/segmentio/ksuid"
)

func httpToWs(t *testing.T, url string) string {
	t.Helper()
	return "ws" + strings.TrimPrefix(url, "http")
}

func newWSServer(t *testing.T, h http.Handler) (*httptest.Server, *websocket.Conn) {
	t.Helper()

	s := httptest.NewServer(h)
	wsURL := httpToWs(t, s.URL)

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	return s, ws
}

func sendMessage(t *testing.T, ws *websocket.Conn, msg Event) {
	t.Helper()

	m, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.WriteMessage(websocket.BinaryMessage, m); err != nil {
		t.Fatalf("%v", err)
	}
}

func receiveWSMessage(t *testing.T, ws *websocket.Conn) Event {
	t.Helper()

	_, m, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	var reply Event
	err = json.Unmarshal(m, &reply)
	if err != nil {
		t.Fatal(err)
	}

	return reply
}

func TestReadRoutineExitWithoutLobby(t *testing.T) {
	lm := &LobbyManager{make(map[ksuid.KSUID]*Lobby)}
	s, ws := newWSServer(t, http.HandlerFunc(lm.ServeWs))

	client := newClient(ws)
	defer s.Close()
	defer ws.Close()

	err := client.readRoutine()

	if err == nil {
		t.Errorf("Read routine should have thrown error when lobby is not set")
	}
}

func TestWriteRoutineExitWithoutLobby(t *testing.T) {
	lm := &LobbyManager{make(map[ksuid.KSUID]*Lobby)}
	s, ws := newWSServer(t, http.HandlerFunc(lm.ServeWs))

	client := newClient(ws)
	defer s.Close()
	defer ws.Close()

	err := client.writeRoutine()

	if err == nil {
		t.Errorf("Read routine should have thrown error when lobby is not set")
	}
}
