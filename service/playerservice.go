package service

import (
	"fmt"
	"github.com/ksanta/word-stallion/dao"
	"github.com/ksanta/word-stallion/model"
	"sync"
)

type PlayerService struct {
	playerDao *dao.PlayerDao
	apiDao    *dao.ApiDao
}

func NewPlayerService(playerDao *dao.PlayerDao, apiDao *dao.ApiDao) *PlayerService {
	return &PlayerService{
		playerDao: playerDao,
		apiDao:    apiDao,
	}
}

func (playerService *PlayerService) SendWelcomeMessageToPlayer(connectionId string, targetScore int) error {
	welcomeMessage := model.MessageToPlayer{
		Welcome: &model.Welcome{TargetScore: targetScore},
	}
	return playerService.apiDao.SendMessageToPlayer(connectionId, welcomeMessage)
}

func (playerService *PlayerService) SendCorrectAnswerToPlayer(connectionId string, correct bool, correctAnswer int) error {
	answerMessage := model.MessageToPlayer{
		PlayerResult: &model.PlayerResult{
			Correct:       correct,
			CorrectAnswer: correctAnswer,
		},
	}
	return playerService.apiDao.SendMessageToPlayer(connectionId, answerMessage)
}

func (playerService *PlayerService) SendRoundSummaryToActivePlayers(gameId string) (model.Players, error) {
	players, err := playerService.playerDao.GetPlayers(gameId)
	if err != nil {
		return nil, fmt.Errorf("error getting players: %w", err)
	}

	roundSummaryMsg := model.MessageToPlayer{
		RoundSummary: &model.RoundSummary{
			PlayerStates: players.PlayerStates(),
		},
	}
	playerService.sendMessageToActivePlayers(players, roundSummaryMsg)
	return players, nil
}

func (playerService *PlayerService) SendAboutToStartToActivePlayers(gameId string, startingInSeconds int) (model.Players, error) {
	players, err := playerService.playerDao.GetPlayers(gameId)
	if err != nil {
		return nil, fmt.Errorf("error getting players: %w", err)
	}

	aboutToStartMsg := &model.MessageToPlayer{
		AboutToStart: &model.AboutToStart{
			Seconds: startingInSeconds,
		},
	}
	playerService.sendMessageToActivePlayers(players, aboutToStartMsg)
	return players, nil
}

func (playerService *PlayerService) SendQuestionToActivePlayers(players model.Players, wordsInThisRound model.Words, correctAnswer int, secondsPerQuestion int) error {
	questionMsg := model.MessageToPlayer{
		PresentQuestion: &model.PresentQuestion{
			WordToGuess:    wordsInThisRound[correctAnswer].Word,
			Definitions:    wordsInThisRound.GetDefinitions(),
			SecondsAllowed: secondsPerQuestion,
		},
	}
	playerService.sendMessageToActivePlayers(players, questionMsg)
	return nil
}

func (playerService *PlayerService) SendGameSummaryToAllActivePlayers(players model.Players) error {
	winner := players.PlayerWithHighestPoints()
	msg := model.MessageToPlayer{
		Summary: &model.Summary{
			Winner: winner.Name,
			Icon:   winner.Icon,
		},
	}
	playerService.sendMessageToActivePlayers(players, msg)
	return nil
}

func (playerService *PlayerService) sendMessageToActivePlayers(players model.Players, message interface{}) {
	waitGroup := sync.WaitGroup{}

	for _, player := range players {
		if player.Active {
			waitGroup.Add(1)
			// Make a copy so goroutine will pick out the correct connection id
			connectionId := player.ConnectionId
			go func() {
				defer waitGroup.Done()
				err := playerService.apiDao.SendMessageToPlayer(connectionId, message)
				if err != nil {
					fmt.Println("Error posting message to player", err)
				}
			}()
		}
	}

	waitGroup.Wait()
}
