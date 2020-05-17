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
	   init {
	     loadWords from S3
	   }

	   DoRound(gameId)

	   if player's max score < target score
	     set game.start_time
	     set all players to waiting
	     pick a question
	     set correct answer on the game
	     send question to all players
	   else
	     send game summary to all active players

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
