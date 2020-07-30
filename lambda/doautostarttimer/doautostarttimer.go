package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ksanta/word-stallion/dao"
	"os"
	"time"
)

var (
	gameDao                 *dao.GameDao
	functionDao             *dao.FunctionDao
	doStartGameFunctionName string
)

func init() {
	gameDao = dao.NewGameDao(os.Getenv("GAMES_TABLE"))
	functionDao = dao.NewFunctionDao()
	doStartGameFunctionName = os.Getenv("DO_START_GAME_FUNCTION_NAME")
}

func handler(gameId string) error {
	game, err := gameDao.GetGame(gameId)
	if err != nil {
		return fmt.Errorf("error getting game: %w\n", err)
	}

	sleepDuration := game.GameStartTime.Sub(time.Now())
	fmt.Println("Sleeping for", sleepDuration.String())
	time.Sleep(sleepDuration)
	return functionDao.InvokeStartGame(doStartGameFunctionName, gameId)
}

func main() {
	lambda.Start(handler)
}
