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
	   player.waiting = false
	   Award points to the player
	   Send msg to player with correct answer

	   if all active players waiting == false
	     sleep 2 seconds
	     invoke "Do Round"
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
