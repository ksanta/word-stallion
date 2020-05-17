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
	CorrectAnswer       int           `json:"correct_answer"`
	GameInProgress      bool          `json:"game_in_progress"`
	WaitingForAnswers   bool          `json:"waiting_for_answers"`
}
