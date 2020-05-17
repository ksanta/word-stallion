package dao

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ksanta/word-stallion/model"
	"time"
)

type GameDao struct {
	service   *dynamodb.DynamoDB
	tableName *string
}

func NewGameDao(tableName string) *GameDao {
	mySession := session.Must(session.NewSession())

	return &GameDao{
		service:   dynamodb.New(mySession),
		tableName: aws.String(tableName),
	}
}

func (gameDao *GameDao) GetPendingGame() (*model.Game, error) {
	// Scan for a pending game
	scanInput := &dynamodb.ScanInput{
		TableName:        gameDao.tableName,
		FilterExpression: aws.String("game_in_progress = :gameInProgress"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":gameInProgress": {
				BOOL: aws.Bool(false),
			},
		},
		ConsistentRead: aws.Bool(true),
	}
	scanOutput, err := gameDao.service.Scan(scanInput)
	if err != nil {
		return nil, err
	}

	if *scanOutput.Count == 0 {
		// If there is no pending game, create one to allow players to group together
		gameId := time.Now().String()
		fmt.Println("Creating a new game:", gameId)

		newGame := &model.Game{
			GameId:              gameId,
			TargetScore:         500,
			OptionsPerQuestion:  3,
			DurationPerQuestion: 10 * time.Second,
			MaxPlayerCount:      5,
			CorrectAnswer:       -1,
			GameInProgress:      false,
			WaitingForAnswers:   false,
		}
		marshalledGame, err := dynamodbattribute.MarshalMap(newGame)
		if err != nil {
			return nil, err
		}
		putItemInput := &dynamodb.PutItemInput{
			TableName: gameDao.tableName,
			Item:      marshalledGame,
		}
		_, err = gameDao.service.PutItem(putItemInput)
		if err != nil {
			return nil, err
		}
		return newGame, nil

	} else if *scanOutput.Count == 1 {
		// Found a pending game so return that
		game := &model.Game{}
		err = dynamodbattribute.UnmarshalMap(scanOutput.Items[0], game)
		if err != nil {
			return nil, err
		}
		fmt.Println("Found pending game:", game.GameId)
		return game, nil
	}
	return nil, errors.New("found more than one pending game")
}

func (gameDao *GameDao) GetGame(gameId string) (*model.Game, error) {
	fmt.Println("Constructing GetItemInput")
	input := &dynamodb.GetItemInput{
		TableName: gameDao.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"game_id": {
				S: aws.String(gameId),
			},
		},
		ConsistentRead: aws.Bool(true),
	}

	fmt.Println("Calling GetItem")
	output, err := gameDao.service.GetItem(input)
	if err != nil {
		return nil, err
	}

	if output.Item == nil {
		return nil, nil
	}

	fmt.Println("Unmarshalling")
	game := &model.Game{}
	err = dynamodbattribute.UnmarshalMap(output.Item, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}
