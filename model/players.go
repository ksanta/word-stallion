package model

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"sync"
)

// todo: change this to NOT a pointer?
type Players []*Player

func (players Players) SendMessageToActivePlayers(message interface{}, event events.APIGatewayWebsocketProxyRequest) error {
	waitGroup := sync.WaitGroup{}

	// todo: create the mgmtService once - move to dao?
	mySession := session.Must(session.NewSession())
	apiMgmtService := apigatewaymanagementapi.New(mySession, &aws.Config{
		Endpoint: aws.String(event.RequestContext.DomainName + "/" + event.RequestContext.Stage),
	})

	for _, player := range players {
		if player.Active {
			waitGroup.Add(1)
			// Make a copy so goroutine will pick out correct connection id
			connectionId := player.ConnectionId
			go func() {
				defer waitGroup.Done()

				marshalledMessage, err := json.Marshal(message)
				if err != nil {
					fmt.Println("Error marshalling welcome message", err)
					return
				}

				postToConnectionInput := &apigatewaymanagementapi.PostToConnectionInput{
					ConnectionId: aws.String(connectionId),
					Data:         marshalledMessage,
				}

				fmt.Println("Posting msg to player", connectionId)
				_, err = apiMgmtService.PostToConnection(postToConnectionInput)
				if err != nil {
					fmt.Println("Error posting welcome message to the player", err)
					return
				}
			}()
		}
	}
	waitGroup.Wait()

	return nil
}

func (players *Players) PlayerStates() []PlayerState {
	playerStates := make([]PlayerState, 0, len(*players))
	for _, p := range *players {
		playerState := PlayerState{
			Name:   p.Name,
			Score:  p.Points,
			Active: p.Active,
			Icon:   p.Icon,
		}
		playerStates = append(playerStates, playerState)
	}

	return playerStates
}

// AllInactive will return true if all the players are inactive
func (players Players) AllInactive() bool {
	for _, p := range players {
		if p.Active {
			return false
		}
	}
	return true
}

func (players Players) NumActivePlayers() int {
	activePlayers := 0

	for _, p := range players {
		if p.Active {
			activePlayers++
		}
	}
	return activePlayers
}

// PlayerWithHighestPoints returns the player with the maximum points. They may not have actually won yet.
func (players Players) PlayerWithHighestPoints() *Player {
	maxScore := -1
	var winner *Player

	for _, p := range players {
		if p.Points > maxScore {
			maxScore = p.Points
			winner = p
		}
	}

	return winner
}
