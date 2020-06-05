// Handles a player message saying that they are ready to play
package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambda2 "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/ksanta/word-stallion/dao"
	"github.com/ksanta/word-stallion/model"
	"github.com/ksanta/word-stallion/service"
	"os"
)

var (
	gameDao                 *dao.GameDao
	playerDao               *dao.PlayerDao
	apiDao                  *dao.ApiDao
	playerService           *service.PlayerService
	lambdaService           *lambda2.Lambda
	doStartGameFunctionName string
)

func init() {
	gameDao = dao.NewGameDao(os.Getenv("GAMES_TABLE"))
	playerDao = dao.NewPlayerDao(os.Getenv("PLAYERS_TABLE"))
	apiDao = dao.NewApiDao(os.Getenv("API_ENDPOINT"))
	playerService = service.NewPlayerService(playerDao, apiDao)

	mySession := session.Must(session.NewSession())
	lambdaService = lambda2.New(mySession)
	doStartGameFunctionName = os.Getenv("DO_START_GAME_FUNCTION_NAME")
}

func handler(event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get a pending game. One will be created if there isn't one yet.
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
	player, err := playerDao.AddNewPlayer(event.RequestContext.ConnectionID,
		game.GameId, newPlayerMessage.Name, newPlayerMessage.Icon)
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

	// Auto-start game if max-players-per-game has been reached
	if len(players) >= game.MaxPlayerCount {
		fmt.Println("Auto-starting game", game.GameId, "after reaching max players")
		err := invokeDoStartGame(game.GameId)
		if err != nil {
			return newErrorResponse("Error invoking start game function", err)
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func invokeDoStartGame(gameId string) error {
	invokeInput := &lambda2.InvokeInput{
		FunctionName:   aws.String(doStartGameFunctionName),
		InvocationType: aws.String(lambda2.InvocationTypeEvent),
		Payload:        []byte("\"" + gameId + "\""),
	}

	_, err := lambdaService.Invoke(invokeInput)
	return err
}

func newErrorResponse(msg string, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
	}, fmt.Errorf("%s: %w", msg, err)
}

func main() {
	lambda.Start(handler)
}
