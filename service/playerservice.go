package service

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
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

func (playerService *PlayerService) SendRoundSummaryToAllActivePlayers(gameId string, event events.APIGatewayWebsocketProxyRequest) (model.Players, error) {
	players, err := playerService.playerDao.GetPlayers(gameId)
	if err != nil {
		return nil, fmt.Errorf("error getting players: %w", err)
	}

	roundSummaryMsg := model.MessageToPlayer{
		RoundSummary: &model.RoundSummary{
			PlayerStates: players.PlayerStates(),
		},
	}

	err = players.SendMessageToActivePlayers(roundSummaryMsg, event)
	if err != nil {
		return nil, fmt.Errorf("error sending msg to players: %w", err)
	}

	return players, nil
}
