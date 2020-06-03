package dao

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ksanta/word-stallion/model"
	"sync"
	"time"
)

type PlayerDao struct {
	service   *dynamodb.DynamoDB
	tableName *string
}

func NewPlayerDao(tableName string) *PlayerDao {
	mySession := session.Must(session.NewSession())

	return &PlayerDao{
		service:   dynamodb.New(mySession),
		tableName: aws.String(tableName),
	}
}

func (playerDao *PlayerDao) AddNewPlayer(connectionId string, gameId string, name string, icon string) error {
	newPlayer := &model.Player{
		ConnectionId:       connectionId,
		GameId:             gameId,
		Active:             true,
		WaitingForResponse: false,
		Name:               name,
		Icon:               icon,
		Points:             0,
		ExpiresAt:          time.Now().Add(10 * time.Minute).Unix(),
	}

	return playerDao.PutPlayer(newPlayer)
}

func (playerDao *PlayerDao) PutPlayer(player *model.Player) error {
	marshalledPlayer, err := dynamodbattribute.MarshalMap(player)
	if err != nil {
		return err
	}

	putItemInput := &dynamodb.PutItemInput{
		Item:      marshalledPlayer,
		TableName: playerDao.tableName,
	}

	_, err = playerDao.service.PutItem(putItemInput)
	return err
}

func (playerDao *PlayerDao) PutPlayers(players model.Players) error {
	waitGroup := sync.WaitGroup{}

	for _, player := range players {
		if player.Active {
			waitGroup.Add(1)
			playerCopy := player
			go func() {
				defer waitGroup.Done()
				err := playerDao.PutPlayer(playerCopy)
				if err != nil {
					fmt.Println("error saving player in PutPlayers", err)
				}
			}()
		}
	}

	waitGroup.Wait()
	// todo: error handling
	return nil
}

func (playerDao *PlayerDao) InactivatePlayer(connectionId string) (*model.Player, error) {
	// Prepare the request
	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: playerDao.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"connection_id": {
				S: aws.String(connectionId),
			},
		},
		UpdateExpression: aws.String("SET #Active = :active"),
		ExpressionAttributeNames: map[string]*string{
			"#Active": aws.String("active"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":active": {
				BOOL: aws.Bool(false),
			},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}

	// Send the update request
	updateItemOutput, err := playerDao.service.UpdateItem(updateItemInput)
	if err != nil {
		return nil, err
	}

	// Unmarshal the DynamoDB map into an object
	player := &model.Player{}
	err = dynamodbattribute.UnmarshalMap(updateItemOutput.Attributes, player)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (playerDao *PlayerDao) GetPlayer(connectionId string) (*model.Player, error) {
	getItemInput := &dynamodb.GetItemInput{
		TableName: playerDao.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"connection_id": {
				S: aws.String(connectionId),
			},
		},
		ConsistentRead: aws.Bool(true),
	}

	output, err := playerDao.service.GetItem(getItemInput)
	if err != nil {
		return nil, fmt.Errorf("playerDao.GetPlayer: %w", err)
	}

	// This is possible if a player doesn't have a record
	if output.Item == nil {
		return nil, nil
	}

	player := &model.Player{}
	err = dynamodbattribute.UnmarshalMap(output.Item, player)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (playerDao *PlayerDao) GetPlayers(gameId string) (model.Players, error) {
	// todo: replace this with a call to an index
	scanInput := &dynamodb.ScanInput{
		TableName:        playerDao.tableName,
		FilterExpression: aws.String("game_id = :gameId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":gameId": {
				S: aws.String(gameId),
			},
		},
		ConsistentRead: aws.Bool(true),
	}
	scanOutput, err := playerDao.service.Scan(scanInput)
	if err != nil {
		return nil, err
	}

	players := make(model.Players, 0, *scanOutput.Count)
	for _, item := range scanOutput.Items {
		player := &model.Player{}
		err = dynamodbattribute.UnmarshalMap(item, player)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

func (playerDao *PlayerDao) DeletePlayer(connectionId string) error {
	deleteItemInput := &dynamodb.DeleteItemInput{
		TableName: playerDao.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"connection_id": {
				S: aws.String(connectionId),
			},
		},
	}

	_, err := playerDao.service.DeleteItem(deleteItemInput)
	return err
}
