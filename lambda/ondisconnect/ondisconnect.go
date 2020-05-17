package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ksanta/word-stallion/dao"
	"github.com/ksanta/word-stallion/service"
	"os"
)

var (
	gameDao       *dao.GameDao
	playerDao     *dao.PlayerDao
	playerService *service.PlayerService
)

func init() {
	gameDao = dao.NewGameDao(os.Getenv("GAMES_TABLE"))
	playerDao = dao.NewPlayerDao(os.Getenv("PLAYERS_TABLE"))
	playerService = service.NewPlayerService(playerDao)
}

func handler(event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionId := event.RequestContext.ConnectionID
	fmt.Println("Disconnect from", connectionId)

	player, err := playerDao.GetPlayer(connectionId)
	if err != nil {
		return newErrorResponse("error getting player", err)
	}

	// Early exit if a person disconnect but doesn't have a Player item
	if player == nil {
		fmt.Println("No player item to update/delete")
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	}

	game, err := gameDao.GetGame(player.GameId)
	if err != nil {
		return newErrorResponse("error getting game", err)
	}

	if !game.GameInProgress {
		fmt.Println("Deleting player", connectionId)
		err = playerDao.DeletePlayer(connectionId)
		if err != nil {
			return newErrorResponse("error deleting player", err)
		}

	} else {
		fmt.Println("Inactivating player", connectionId)
		_, err = playerDao.InactivatePlayer(connectionId)
		if err != nil {
			return newErrorResponse("error inactivating player", err)
		}
	}

	_, err = playerService.SendRoundSummaryToAllActivePlayers(game.GameId, event)
	if err != nil {
		return newErrorResponse("error ending round update to all players", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func newErrorResponse(msg string, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
	}, fmt.Errorf("%s: %w", msg, err)
}

func main() {
	lambda.Start(handler)
}
