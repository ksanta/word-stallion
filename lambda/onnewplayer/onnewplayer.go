// Handles a player message saying that they are ready to play
package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ksanta/word-stallion/dao"
	"github.com/ksanta/word-stallion/model"
	"github.com/ksanta/word-stallion/service"
	"os"
	"time"
)

var (
	gameDao                 *dao.GameDao
	playerDao               *dao.PlayerDao
	apiDao                  *dao.ApiDao
	playerService           *service.PlayerService
	functionDao             *dao.FunctionDao
	doStartGameFunctionName string
	doAutostartTimerName    string
)

func init() {
	gameDao = dao.NewGameDao(os.Getenv("GAMES_TABLE"))
	playerDao = dao.NewPlayerDao(os.Getenv("PLAYERS_TABLE"))
	apiDao = dao.NewApiDao(os.Getenv("API_ENDPOINT"))
	playerService = service.NewPlayerService(playerDao, apiDao)
	functionDao = dao.NewFunctionDao()

	doStartGameFunctionName = os.Getenv("DO_START_GAME_FUNCTION_NAME")
	doAutostartTimerName = os.Getenv("DO_AUTOSTART_TIMER_FUNCTION_NAME")
}

func handler(event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get a pending game. One will be created if there isn't one yet.
	// todo: move logic into here
	game, err := gameDao.GetPendingGame()
	if err != nil {
		return newErrorResponse("Failed to get Game item", err)
	}

	// Extract the player's message from the event
	fmt.Println("Received new player msg:", event.Body)
	var playerMessage model.MessageFromPlayer
	err = json.Unmarshal([]byte(event.Body), &playerMessage)
	if err != nil {
		return newErrorResponse("Error unmarshalling JSON body", err)
	}
	newPlayerMessage := playerMessage.NewPlayer

	// Create a new Player item in Dynamo
	fmt.Println("Saving new player:", event.RequestContext.ConnectionID)
	millisSinceGameCreated := time.Since(game.CreatedAt).Milliseconds()
	player, err := playerDao.AddNewPlayer(event.RequestContext.ConnectionID,
		game.GameId, millisSinceGameCreated, newPlayerMessage.Name, newPlayerMessage.Icon)
	if err != nil {
		return newErrorResponse("Error saving new player", err)
	}

	// Send a welcome message to the player
	err = playerService.SendWelcomeMessageToPlayer(*player, game.TargetScore)
	if err != nil {
		return newErrorResponse("Error posting welcome message to the player", err)
	}

	// Send a "round summary" message to all active players
	players, err := playerService.SendRoundSummaryToActivePlayers(game.GameId)
	if err != nil {
		return newErrorResponse("Error sending a message to all players", err)
	}

	// If this is the first player, invoke the auto start timer
	// todo: alternately, invoke this when a game is created
	if len(players) == 1 {
		err := functionDao.InvokeAutostartTimer(doAutostartTimerName, game.GameId)
		if err != nil {
			return newErrorResponse("Error invoking autostart timer function", err)
		}
	}

	// Auto-start game if max-players-per-game has been reached
	if len(players) >= game.MaxPlayerCount {
		fmt.Println("Auto-starting game", game.GameId, "after reaching max players")
		err := functionDao.InvokeStartGame(doStartGameFunctionName, game.GameId)
		if err != nil {
			return newErrorResponse("Error invoking start game function", err)
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
