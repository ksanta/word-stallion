package model

// todo: change this to NOT a pointer?
type Players []*Player

func (players *Players) PlayerStates() []PlayerState {
	playerStates := make([]PlayerState, 0, len(*players))
	for _, p := range *players {
		playerState := PlayerState{
			Id:     p.ConnectionId,
			Name:   p.Name,
			Score:  p.Points,
			Active: p.Active,
			Icon:   p.Icon,
		}
		playerStates = append(playerStates, playerState)
	}

	return playerStates
}

// AllInactive will return true if all the players are inactive
func (players Players) AllInactive() bool {
	for _, p := range players {
		if p.Active {
			return false
		}
	}
	return true
}

func (players Players) NumActivePlayers() int {
	activePlayers := 0

	for _, p := range players {
		if p.Active {
			activePlayers++
		}
	}
	return activePlayers
}

// PlayerWithHighestPoints returns the player with the maximum points. They may not have actually won yet.
func (players Players) PlayerWithHighestPoints() *Player {
	maxScore := -1
	var winner *Player

	for _, p := range players {
		if p.Points > maxScore {
			maxScore = p.Points
			winner = p
		}
	}

	return winner
}

func (players Players) AllActivePlayersResponded() bool {
	for _, player := range players {
		if player.Active && !player.Responded {
			return false
		}
	}
	return true
}

func (players Players) SetActivesToNotResponded() {
	for _, p := range players {
		if p.Active {
			p.Responded = false
		}
	}
}
