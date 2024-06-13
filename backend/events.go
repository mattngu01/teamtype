package main

const (
	LobbyInfo = "LobbyInfo"
)

// Representation of data passed from backend / frontend
type Event struct {
	Type string `json:"type"`
	Data map[string]interface{} `json:"data"`
}
