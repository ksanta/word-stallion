package model

import (
	"time"
)

type Game struct {
	GameId              string        `json:"game_id"`
	TargetScore         int           `json:"target_score"`
	OptionsPerQuestion  int           `json:"options_per_question"`
	DurationPerQuestion time.Duration `json:"duration_per_question"`
	MaxPlayerCount      int           `json:"max_player_count"`
	GameInProgress      bool          `json:"game_in_progress"`
	CorrectAnswer       int           `json:"correct_answer"`
	StartTime           time.Time     `json:"start_time"`
	ExpiresAt           int64         `json:"expires_at"`
}
