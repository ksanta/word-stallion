package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ksanta/word-stallion/dao"
	"os"
	"time"
)

var (
	functionDao             *dao.FunctionDao
	doStartGameFunctionName string
)

func init() {
	functionDao = dao.NewFunctionDao()
	doStartGameFunctionName = os.Getenv("DO_START_GAME_FUNCTION_NAME")
}

func handler(gameId string) error {
	time.Sleep(30 * time.Second)
	return functionDao.InvokeStartGame(doStartGameFunctionName, gameId)
}

func main() {
	lambda.Start(handler)
}
