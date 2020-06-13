package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ksanta/word-stallion/dao"
	"github.com/ksanta/word-stallion/model"
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
	apiDao := dao.NewApiDao(os.Getenv("API_ENDPOINT"))
	playerService = service.NewPlayerService(playerDao, apiDao)
}

func handler(event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionId := event.RequestContext.ConnectionID
	fmt.Println("Getting player")
	player, err := playerDao.GetPlayer(connectionId)
	if err != nil {
		return newErrorResponse("error getting player", err)
	}

	// Ignore if a person disconnected but doesn't have a Player item
	if player == nil {
		fmt.Println("Disconnect from unregistered player")
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	} else {
		fmt.Println(player.Name, "disconnected")
	}

	fmt.Println("Getting game")
	game, err := gameDao.GetGame(player.GameId)
	if err != nil {
		return newErrorResponse("error getting game", err)
	}

	// Ignore disconnect if the game is finished
	if game.GameState == model.Finished {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	}

	if game.GameState == model.Pending {
		fmt.Println("Deleting", player.Name)
		err = playerDao.DeletePlayer(connectionId)
		if err != nil {
			return newErrorResponse("error deleting player", err)
		}

	} else if game.GameState == model.InProgress {
		fmt.Println("Inactivating", player.Name)
		_, err = playerDao.InactivatePlayer(connectionId)
		if err != nil {
			return newErrorResponse("error inactivating player", err)
		}
	}

	players, err := playerService.SendRoundSummaryToActivePlayers(game.GameId)
	if err != nil {
		return newErrorResponse("error ending round update to all players", err)
	}

	// Delete the pending game if all players are inactive
	if game.GameState == model.Pending && players.AllInactive() {
		err := gameDao.DeleteGame(game)
		if err != nil {
			return newErrorResponse("error deleting pending game", err)
		}
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
