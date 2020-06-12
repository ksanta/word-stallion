package model

import (
	"time"
)

type Game struct {
	GameId             string    `json:"game_id"`
	TargetScore        int       `json:"target_score"`
	OptionsPerQuestion int       `json:"options_per_question"`
	SecondsPerQuestion int       `json:"duration_per_question"`
	MaxPlayerCount     int       `json:"max_player_count"`
	GameState          GameState `json:"game_state"`
	CorrectAnswer      int       `json:"correct_answer"`
	RoundStartTime     time.Time `json:"round_start_time"`
	CreatedAt          time.Time `json:"created_at"`
	ExpiresAt          int64     `json:"expires_at"`
}

type GameState string

const (
	Pending    = GameState("PENDING")
	InProgress = GameState("IN_PROGRESS")
	Finished   = GameState("FINISHED")
)

func (game *Game) CalculatePoints(submittedAnswer int, timeReceived time.Time) int {
	elapsedDuration := timeReceived.Sub(game.RoundStartTime)
	durationPerQuestion := time.Duration(game.SecondsPerQuestion) * time.Second

	// Player took longer than allowed time - no points!
	if elapsedDuration > durationPerQuestion {
		return 0
	}

	// Points for correct answer
	correctPoints := 0
	if submittedAnswer == game.CorrectAnswer {
		correctPoints += 100
	}

	// Points for answering quickly
	timePoints := int(50 * (durationPerQuestion - elapsedDuration) / durationPerQuestion)
	if timePoints < 0 {
		timePoints = 0
	}

	return correctPoints + timePoints
}
