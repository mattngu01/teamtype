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

func TestRespondJoinLobbyWithLobbyInfo(t *testing.T) {
	lm := &LobbyManager{make(map[ksuid.KSUID]*Lobby)}
	server, frontendWebsocket := newWSServer(t, http.HandlerFunc(lm.ServeWs))
	defer server.Close()
	defer frontendWebsocket.Close()

	joinLobbyRequest := Event{
		Type: JoinLobby,
		Data: JoinLobbyData{
			Username: "test_user",
		},
	}
	jsonJoinLobbyRequest, _ := json.Marshal(joinLobbyRequest)
	// want to send request to server, for join lobby
	if err := frontendWebsocket.WriteMessage(websocket.TextMessage, jsonJoinLobbyRequest); err != nil {
		t.Fatalf("Unable to send message through websocket %s", err)
	}
	_, resp, err := frontendWebsocket.ReadMessage()
	if err != nil {
		t.Fatalf("Could not read message from websocket %s", err)
	}

	lobbyInfoResponse := Event{}
	err = json.Unmarshal(resp, &lobbyInfoResponse)

	if err != nil {
		t.Fatalf("Could not unmarshal server response %s", err)
	}

	if lobbyInfoResponse.Type != LobbyInfo {
		t.Fatalf("Response was not type Lobby Info")
	}

	// if lobbyInfoResponse.Data.Players[0] != "test_user" {
	// 	t.Fatalf("Response did not include client in lobby")
	// }
}
