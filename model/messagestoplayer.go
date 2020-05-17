package model

type MessageToPlayer struct {
	//PlayerDetailsReq *PlayerDetailsReq `json:",omitempty"`
	Welcome *Welcome `json:",omitempty"`
	//AboutToStart     *AboutToStart     `json:",omitempty"`
	//PresentQuestion  *PresentQuestion  `json:",omitempty"`
	//PlayerResult     *PlayerResult     `json:",omitempty"`
	RoundSummary *RoundSummary `json:",omitempty"`
	//Summary          *Summary          `json:",omitempty"`
	//Error            *GameError        `json:",omitempty"`
}

// Welcome tells the client to display an intro to the player
type Welcome struct {
	TargetScore int
}

type RoundSummary struct {
	PlayerStates []PlayerState
}

type PlayerState struct {
	Name   string
	Icon   string
	Score  int
	Active bool
}
