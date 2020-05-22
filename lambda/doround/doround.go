package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ksanta/word-stallion/dao"
	"github.com/ksanta/word-stallion/model"
	"github.com/ksanta/word-stallion/service"
	"os"
)

var (
	gameDao       *dao.GameDao
	playerDao     *dao.PlayerDao
	playerService *service.PlayerService
	words         model.Words
)

func init() {
	gameDao = dao.NewGameDao(os.Getenv("GAMES_TABLE"))
	playerDao = dao.NewPlayerDao(os.Getenv("PLAYERS_TABLE"))
	playerService = service.NewPlayerService(playerDao)

	bucketName := os.Getenv("WORDS_BUCKET")
	wordsDao := dao.NewWordsDao(bucketName)
	words2, err := wordsDao.GetWords()
	if err != nil {
		fmt.Println("error loading words:", err)
		return
	}
	words = words2
}

func handler(gameId string) error {
	// Fetch the info we need
	players, err := playerDao.GetPlayers(gameId)
	if err != nil {
		return fmt.Errorf("error getting players: %w\n", err)
	}
	game, err := gameDao.GetGame(gameId)
	if err != nil {
		return fmt.Errorf("error getting game: %w\n", err)
	}

	if players.PlayerWithHighestPoints().Points < game.TargetScore {
		fmt.Println("todo: send question to all players")
		/*
		 prepare question and answer

		 set game.start_time
		 set game.correct_answer
		 save game

		 set players.waiting
		 save players

		 send question to all players
		*/

	} else {
		// 	send game summary to all active players
		fmt.Println("todo: send game summary to all players")
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
