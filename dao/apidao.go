package dao

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/ksanta/word-stallion/model"
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

func (apiDao *ApiDao) SendMessageToPlayer(player model.Player, message interface{}, msgType string) error {
	fmt.Printf("Sending %s to %s\n", msgType, player.Name)

	marshalledMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	postToConnectionInput := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(player.ConnectionId),
		Data:         marshalledMessage,
	}

	_, err = apiDao.service.PostToConnection(postToConnectionInput)
	return err
}
