package model

// Player represents a player that started playing a game. The JSON metadata is for converting
// this struct into a DynamoDB item.
type Player struct {
	ConnectionId string `json:"connection_id"`
	// The id of the game this player is a part of
	GameId string `json:"game_id"`
	// Whether the player has an active connection. Connection could go dead mid-game
	// and the show must go on!
	Active bool `json:"active"`
	// Milliseconds since the game was created when the player joined
	MillisSinceGameCreatedWhenJoined int64 `json:"millis_since_game_created_when_joined"`
	// This player has responded to the question
	Responded bool `json:"responded"`
	// Name of this player
	Name string `json:"name"`
	// Client-specific icon to represent the player
	Icon string `json:"icon"`
	// Points for this player
	Points int `json:"points"`
	// Time this record will expire
	ExpiresAt int64 `json:"expires_at"`
}
