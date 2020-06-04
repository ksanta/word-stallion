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
	"time"
)

var (
	gameDao             *dao.GameDao
	playerDao           *dao.PlayerDao
	playerService       *service.PlayerService
	lambdaService       *lambda2.Lambda
	doRoundFunctionName string
)

func init() {
	gameDao = dao.NewGameDao(os.Getenv("GAMES_TABLE"))
	playerDao = dao.NewPlayerDao(os.Getenv("PLAYERS_TABLE"))
	apiDao := dao.NewApiDao(os.Getenv("API_ENDPOINT"))
	playerService = service.NewPlayerService(playerDao, apiDao)

	mySession := session.Must(session.NewSession())
	lambdaService = lambda2.New(mySession)
	doRoundFunctionName = os.Getenv("DO_ROUND_FUNCTION_NAME")
}

func handler(event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get the information we need
	player, err := playerDao.GetPlayer(event.RequestContext.ConnectionID)
	if err != nil {
		return newErrorResponse("error fetching player", err)
	}
	// Early exit if the player has already submitted their response
	if player.WaitingForResponse == false {
		fmt.Println("Player already responded - ignoring")
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	}
	player.WaitingForResponse = false

	game, err := gameDao.GetGame(player.GameId)
	if err != nil {
		return newErrorResponse("error fetching game", err)
	}

	// Extract player response from the request
	fmt.Println("Received new player msg:", event.Body)
	playerMessage := model.MessageFromPlayer{}
	err = json.Unmarshal([]byte(event.Body), &playerMessage)
	if err != nil {
		return newErrorResponse("error unmarshalling JSON body", err)
	}

	// Award points to the player
	playerResponse := playerMessage.PlayerResponse.Response
	fmt.Println("Player response is", playerResponse)
	timeReceived := time.Now()
	fmt.Println("Time received is", timeReceived)
	pointsForRound := game.CalculatePoints(playerResponse, timeReceived)
	fmt.Println("Player awarded", pointsForRound, "points")
	player.Points += pointsForRound

	// Save player's updated attributes
	err = playerDao.PutPlayer(player)
	if err != nil {
		return newErrorResponse("error saving player", err)
	}

	// Send the correct answer to the player
	err = playerService.SendCorrectAnswerToPlayer(player.ConnectionId, playerResponse == game.CorrectAnswer, game.CorrectAnswer)
	if err != nil {
		return newErrorResponse("error sending correct answer to player", err)
	}

	// If all players have responded, send round update, and do another round or finish the game
	players, err := playerDao.GetPlayers(game.GameId)
	if err != nil {
		return newErrorResponse("error fetching players", err)
	}
	if players.AllActivePlayersResponded() {
		_, err := playerService.SendRoundSummaryToActivePlayers(game.GameId)
		if err != nil {
			return newErrorResponse("error sending round summary", err)
		}
		if players.PlayerWithHighestPoints().Points < game.TargetScore {
			// Do another round if the target score is not yet reached
			time.Sleep(2 * time.Second)
			err = invokeDoRound(game.GameId)
			if err != nil {
				return newErrorResponse("error invoking DoRound", err)
			}
		} else {
			// Finish the game if the target score is reached
			err = playerService.SendGameSummaryToAllActivePlayers(players)
			if err != nil {
				return newErrorResponse("error sending game summary to players", err)
			}
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func invokeDoRound(gameId string) error {
	invokeInput := &lambda2.InvokeInput{
		FunctionName:   aws.String(doRoundFunctionName),
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
