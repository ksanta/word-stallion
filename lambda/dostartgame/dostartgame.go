package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ksanta/word-stallion/dao"
	"github.com/ksanta/word-stallion/service"
	"os"
	"time"
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

func handler(gameId string) error {
	game, err := gameDao.GetGame(gameId)
	if err != nil {
		return fmt.Errorf("error getting game: %w", err)
	}

	// Ignore multiple requests to start a game if it's already in progress
	if game.GameInProgress {
		return nil
	}

	// Update game to in progress
	_, err = gameDao.UpdateToInProgress(gameId)
	if err != nil {
		return fmt.Errorf("error updating game to in progress: %w", err)
	}

	// Send "about to start" message to all active players
	const startingInSeconds = 5
	_, err = playerService.SendAboutToStartToAllActivePlayers(game.GameId, game.Endpoint, startingInSeconds)
	if err != nil {
		return fmt.Errorf("error sending msg to all players: %w", err)
	}

	// Sleep for a bit
	time.Sleep(startingInSeconds * time.Second)

	// Asynchronously invoke DoRound function
	fmt.Println("todo: invoke DoRound function")

	return nil
}

func main() {
	lambda.Start(handler)
}
