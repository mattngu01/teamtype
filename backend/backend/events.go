package backend

const (
	LobbyInfo = "LobbyInfo"
	JoinLobby = "JoinLobby"
)

// Representation of data passed from backend / frontend
type Event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type LobbyInfoData struct {
	LobbyId string   `json:"lobbyId"`
	Players []string `json:"players"`
}

type JoinLobbyData struct {
	Username string `json:"username"`
}
