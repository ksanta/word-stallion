package dao

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

type ApiDao struct {
	service *apigatewaymanagementapi.ApiGatewayManagementApi
}

func NewApiDao(endpoint string) *ApiDao {
	mySession := session.Must(session.NewSession())
	service := apigatewaymanagementapi.New(mySession, &aws.Config{
		Endpoint: aws.String(endpoint),
	})

	return &ApiDao{
		service: service,
	}
}

func (apiDao *ApiDao) SendMessageToPlayer(connectionId string, message interface{}) error {
	fmt.Println("Sending msg to player", connectionId)

	marshalledMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	postToConnectionInput := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionId),
		Data:         marshalledMessage,
	}

	_, err = apiDao.service.PostToConnection(postToConnectionInput)
	return err
}
