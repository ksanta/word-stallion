package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambda2 "github.com/aws/aws-sdk-go/service/lambda"
	"os"
	"time"
)

var (
	lambdaService           *lambda2.Lambda
	doStartGameFunctionName string
)

func init() {
	mySession := session.Must(session.NewSession())
	lambdaService = lambda2.New(mySession)
	doStartGameFunctionName = os.Getenv("DO_START_GAME_FUNCTION_NAME")
}

func handler(gameId string) error {
	time.Sleep(30 * time.Second)
	return invokeDoStartGame(gameId)
}

// todo: move this to a DAO
func invokeDoStartGame(gameId string) error {
	invokeInput := &lambda2.InvokeInput{
		FunctionName:   aws.String(doStartGameFunctionName),
		InvocationType: aws.String(lambda2.InvocationTypeEvent),
		Payload:        []byte("\"" + gameId + "\""),
	}

	_, err := lambdaService.Invoke(invokeInput)
	return err
}

func main() {
	lambda.Start(handler)
}
