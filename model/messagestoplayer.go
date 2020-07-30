package model

type MessageToPlayer struct {
	//PlayerDetailsReq *PlayerDetailsReq `json:",omitempty"`
	Welcome         *Welcome         `json:",omitempty"`
	AboutToStart    *AboutToStart    `json:",omitempty"`
	PresentQuestion *PresentQuestion `json:",omitempty"`
	PlayerResult    *PlayerResult    `json:",omitempty"`
	RoundSummary    *RoundSummary    `json:",omitempty"`
	Summary         *Summary         `json:",omitempty"`
	//Error            *GameError        `json:",omitempty"`
}

// Welcome is sent to a player as they are waiting for the game to start
type Welcome struct {
	SecondsTillStart int
	TargetScore      int
}

// AboutToStart tells all players that the game will start in X seconds
type AboutToStart struct {
	Seconds int
}

// PresentQuestion is the question sent to each player
type PresentQuestion struct {
	WordToGuess    string
	Definitions    []string
	SecondsAllowed int
}

// PlayerResult is sent to the player telling them their result of the round
type PlayerResult struct {
	Correct       bool // todo: drop this field
	CorrectAnswer int
}

// RoundSummary is sent to each active player at the end of each round
type RoundSummary struct {
	PlayerStates []PlayerState
}

// Summary is sent to the client at the end telling the player the final result
type Summary struct {
	Winner string
	Icon   string
}

// PlayerState is a summary of player info as part of the round summary
type PlayerState struct {
	Id     string
	Name   string
	Icon   string
	Score  int
	Active bool
}
