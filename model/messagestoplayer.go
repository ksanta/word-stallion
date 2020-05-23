package model

type MessageToPlayer struct {
	//PlayerDetailsReq *PlayerDetailsReq `json:",omitempty"`
	Welcome         *Welcome         `json:",omitempty"`
	AboutToStart    *AboutToStart    `json:",omitempty"`
	PresentQuestion *PresentQuestion `json:",omitempty"`
	//PlayerResult     *PlayerResult     `json:",omitempty"`
	RoundSummary *RoundSummary `json:",omitempty"`
	//Summary          *Summary          `json:",omitempty"`
	//Error            *GameError        `json:",omitempty"`
}

// Welcome tells the client to display an intro to the player
type Welcome struct {
	TargetScore int
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

// RoundSummary is sent to each active player at the end of each round
type RoundSummary struct {
	PlayerStates []PlayerState
}

// PlayerState is a summary of player info as part of the round summary
type PlayerState struct {
	Name   string
	Icon   string
	Score  int
	Active bool
}
