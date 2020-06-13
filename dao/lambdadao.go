package dao

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type FunctionDao struct {
	lambdaService *lambda.Lambda
}

func NewFunctionDao() *FunctionDao {
	mySession := session.Must(session.NewSession())
	return &FunctionDao{
		lambdaService: lambda.New(mySession),
	}
}

func (functionDao *FunctionDao) InvokeAutostartTimer(functionName string, gameId string) error {
	return functionDao.invokeGameFunction(functionName, gameId)
}

func (functionDao *FunctionDao) InvokeStartGame(functionName string, gameId string) error {
	return functionDao.invokeGameFunction(functionName, gameId)
}

func (functionDao *FunctionDao) InvokeDoRound(functionName string, gameId string) error {
	return functionDao.invokeGameFunction(functionName, gameId)
}

func (functionDao *FunctionDao) invokeGameFunction(functionName string, gameId string) error {
	invokeInput := &lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		InvocationType: aws.String(lambda.InvocationTypeEvent),
		Payload:        []byte("\"" + gameId + "\""),
	}

	_, err := functionDao.lambdaService.Invoke(invokeInput)
	return err
}
