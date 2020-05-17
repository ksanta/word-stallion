package model

import "time"

// Player represents a player that started playing a game. The JSON metadata is for converting
// this struct into a DynamoDB item.
type Player struct {
	ConnectionId string `json:"connection_id"`
	// The id of the game this player is a part of
	GameId string `json:"game_id"`
	// Whether the player has an active connection. Connection could go dead mid-game
	// and the show must go on!
	Active bool `json:"active"`
	// WaitingFOrResponse means the player has been sent the question and hasn't received response yet
	WaitingForResponse bool `json:"waiting_for_response"`
	// Name of this player
	Name string `json:"name"`
	// Client-specific icon to represent the player
	Icon string `json:"icon"`
	// Points for this player
	Points int `json:"points"`
	// Time tracks when a player started to answer a question
	StartTime time.Time `json:"start_time"`
}
