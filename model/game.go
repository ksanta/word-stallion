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
	GameInProgress     bool      `json:"game_in_progress"`
	CorrectAnswer      int       `json:"correct_answer"`
	StartTime          time.Time `json:"start_time"`
	ExpiresAt          int64     `json:"expires_at"`
}

func (game *Game) CalculatePoints(submittedAnswer int, timeReceived time.Time) int {
	elapsedDuration := timeReceived.Sub(game.StartTime)
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