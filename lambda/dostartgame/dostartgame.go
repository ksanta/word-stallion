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
	/*
	   StartGame(gameId)

	   if game is in progress
	     return

	   set game in progress to true
	   alert players the game will begin
	   sleep 5 seconds
	   invoke DoRound

	*/
}

func newErrorResponse(msg string, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
	}, fmt.Errorf("%s: %w", msg, err)
}

func main() {
	lambda.Start(handler)
}
