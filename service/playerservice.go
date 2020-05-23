package service

import (
	"fmt"
	"github.com/ksanta/word-stallion/dao"
	"github.com/ksanta/word-stallion/model"
)

type PlayerService struct {
	playerDao *dao.PlayerDao
}

func NewPlayerService(playerDao *dao.PlayerDao) *PlayerService {
	return &PlayerService{
		playerDao: playerDao,
	}
}

func (playerService *PlayerService) SendRoundSummaryToAllActivePlayers(gameId string, endpoint string) (model.Players, error) {
	players, err := playerService.playerDao.GetPlayers(gameId)
	if err != nil {
		return nil, fmt.Errorf("error getting players: %w", err)
	}

	roundSummaryMsg := model.MessageToPlayer{
		RoundSummary: &model.RoundSummary{
			PlayerStates: players.PlayerStates(),
		},
	}

	err = players.SendMessageToActivePlayers(roundSummaryMsg, endpoint)
	if err != nil {
		return nil, fmt.Errorf("error sending msg to players: %w", err)
	}

	return players, nil
}

func (playerService *PlayerService) SendAboutToStartToAllActivePlayers(gameId string, endpoint string, startingInSeconds int) (model.Players, error) {
	players, err := playerService.playerDao.GetPlayers(gameId)
	if err != nil {
		return nil, fmt.Errorf("error getting players: %w", err)
	}

	aboutToStartMsg := &model.MessageToPlayer{
		AboutToStart: &model.AboutToStart{
			Seconds: startingInSeconds,
		},
	}

	err = players.SendMessageToActivePlayers(aboutToStartMsg, endpoint)
	if err != nil {
		return nil, fmt.Errorf("error sending msg to all players: %w", err)
	}

	return players, nil
}

func (playerService *PlayerService) SendQuestionToAllActivePlayers(players model.Players, endpoint string, wordsInThisRound model.Words, correctAnswer int, secondsPerQuestion int) error {
	questionMsg := model.MessageToPlayer{
		PresentQuestion: &model.PresentQuestion{
			WordToGuess:    wordsInThisRound[correctAnswer].Word,
			Definitions:    wordsInThisRound.GetDefinitions(),
			SecondsAllowed: secondsPerQuestion,
		},
	}
	err := players.SendMessageToActivePlayers(questionMsg, endpoint)
	if err != nil {
		return fmt.Errorf("error sending question to all players: %w", err)
	}
	return nil
}
