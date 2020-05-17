package model

type MessageFromPlayer struct {
	// The action to take
	MessageType    string
	NewPlayer      *NewPlayer      `json:",omitempty"`
	PlayerResponse *PlayerResponse `json:",omitempty"`
}

// NewPlayer is sent from the player when they are ready to start playing
type NewPlayer struct {
	Name string
	Icon string
}

// PlayerResponse is the response from the player
type PlayerResponse struct {
	Response int
}
