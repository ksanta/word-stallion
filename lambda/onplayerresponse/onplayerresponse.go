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
	gameDao             *dao.GameDao
	playerDao           *dao.PlayerDao
	playerService       *service.PlayerService
	functionDao         *dao.FunctionDao
	doRoundFunctionName string
)

func init() {
	gameDao = dao.NewGameDao(os.Getenv("GAMES_TABLE"))
	playerDao = dao.NewPlayerDao(os.Getenv("PLAYERS_TABLE"))
	apiDao := dao.NewApiDao(os.Getenv("API_ENDPOINT"))
	playerService = service.NewPlayerService(playerDao, apiDao)

	functionDao = dao.NewFunctionDao()
	doRoundFunctionName = os.Getenv("DO_ROUND_FUNCTION_NAME")
}

func handler(event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Getting player")
	player, err := playerDao.GetPlayer(event.RequestContext.ConnectionID)
	if err != nil {
		return newErrorResponse("error fetching player", err)
	}

	// Early exit if the player has already submitted their response
	if player.Responded == true {
		fmt.Println("Player already responded - ignoring")
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	}
	player.Responded = true

	fmt.Println("Getting game")
	game, err := gameDao.GetGame(player.GameId)
	if err != nil {
		return newErrorResponse("error fetching game", err)
	}

	// Extract player response from the request
	playerMessage := model.MessageFromPlayer{}
	err = json.Unmarshal([]byte(event.Body), &playerMessage)
	if err != nil {
		return newErrorResponse("error unmarshalling JSON body", err)
	}

	// Award points to the player
	playerResponse := playerMessage.PlayerResponse.Response
	fmt.Printf("%s responded with %d\n", player.Name, playerResponse)
	pointsForRound := game.CalculatePoints(playerResponse, time.Now())
	fmt.Printf("%s awarded %d points\n", player.Name, pointsForRound)
	player.Points += pointsForRound

	// Save player's updated attributes
	fmt.Println("Saving player")
	err = playerDao.PutPlayer(player)
	if err != nil {
		return newErrorResponse("error saving player", err)
	}

	// Send the correct answer to the player
	go func() {
		err = playerService.SendCorrectAnswerToPlayer(*player, playerResponse == game.CorrectAnswer, game.CorrectAnswer)
		if err != nil {
			fmt.Println("error sending correct answer to player: %w", err)
		}
	}()

	// If all players have responded, send round update, and do another round or finish the game
	fmt.Println("Getting players")
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
			fmt.Println("Sleeping two seconds")
			time.Sleep(2 * time.Second)
			fmt.Print("Invoking DoRound")
			err = functionDao.InvokeDoRound(doRoundFunctionName, game.GameId)
			if err != nil {
				return newErrorResponse("error invoking DoRound", err)
			}
		} else {
			// Game is finished - send winner to all players
			err = playerService.SendGameSummaryToAllActivePlayers(players)
			if err != nil {
				return newErrorResponse("error sending game summary to players", err)
			}

			// Update game state as finished
			fmt.Println("Updating game as finished")
			game.GameState = model.Finished
			err := gameDao.PutGame(game)
			if err != nil {
				return newErrorResponse("error saving game", err)
			}
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
