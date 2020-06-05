package main

import (
	"fmt"
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

func handler(gameId string) error {
	game, err := gameDao.GetGame(gameId)
	if err != nil {
		return fmt.Errorf("error getting game: %w", err)
	}

	// Ignore multiple requests to start a game if it's already in progress
	if game.GameState == model.InProgress {
		return nil
	}

	// Update game to in progress
	game.GameState = model.InProgress
	game.ExpiresAt = time.Now().Add(10 * time.Minute).Unix()
	err = gameDao.PutGame(game)
	if err != nil {
		return fmt.Errorf("error updating game to in progress: %w", err)
	}

	// todo: Set the expiresAt attribute on all the players

	// Send "about to start" message to all active players
	const startingInSeconds = 5
	_, err = playerService.SendAboutToStartToActivePlayers(game.GameId, startingInSeconds)
	if err != nil {
		return fmt.Errorf("error sending msg to all players: %w", err)
	}

	// Sleep for a bit
	fmt.Println("Sleeping for", startingInSeconds)
	time.Sleep(startingInSeconds * time.Second)

	// Asynchronously invoke DoRound function
	fmt.Println("Invoking function", doRoundFunctionName)
	err = invokeDoRound(gameId)
	if err != nil {
		return fmt.Errorf("error invoking DoRound function: %w", err)
	}

	return nil
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

func main() {
	lambda.Start(handler)
}
