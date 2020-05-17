// Handles a player message saying that they are ready to play
package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
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
	playerService = service.NewPlayerService(playerDao)
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
	err = playerDao.AddNewPlayer(event.RequestContext.ConnectionID,
		game.GameId, newPlayerMessage.Name, newPlayerMessage.Icon)
	if err != nil {
		return newErrorResponse("Error saving new player", err)
	}

	// Send a welcome message to the player
	mySession := session.Must(session.NewSession())
	apiMgmtService := apigatewaymanagementapi.New(mySession, &aws.Config{
		Endpoint: aws.String(event.RequestContext.DomainName + "/" + event.RequestContext.Stage),
	})

	welcomeMessage := model.MessageToPlayer{
		Welcome: &model.Welcome{TargetScore: game.TargetScore},
	}
	marshalledMessage, err := json.Marshal(&welcomeMessage)
	if err != nil {
		return newErrorResponse("Error marshalling welcome message", err)
	}

	fmt.Println("Sending welcome to player:", string(marshalledMessage))
	postToConnectionInput := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(event.RequestContext.ConnectionID),
		Data:         marshalledMessage,
	}

	_, err = apiMgmtService.PostToConnection(postToConnectionInput)
	if err != nil {
		return newErrorResponse("Error posting welcome message to the player", err)
	}

	// Send a "round summary" message to all active players
	players, err := playerService.SendRoundSummaryToAllActivePlayers(game.GameId, event)
	if err != nil {
		return newErrorResponse("Error sending a message to all players", err)
	}

	// Auto-start game if max-players-per-game has been reached
	if len(players) >= game.MaxPlayerCount {
		fmt.Println("TODO: Auto-starting game", game.GameId)
		// todo: invoke doStartGame
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
