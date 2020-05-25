package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ksanta/word-stallion/dao"
	"github.com/ksanta/word-stallion/model"
	"github.com/ksanta/word-stallion/service"
	"os"
	"time"
)

var (
	gameDao       *dao.GameDao
	playerDao     *dao.PlayerDao
	playerService *service.PlayerService
	wordsByType   map[string]model.Words
)

func init() {
	gameDao = dao.NewGameDao(os.Getenv("GAMES_TABLE"))
	playerDao = dao.NewPlayerDao(os.Getenv("PLAYERS_TABLE"))
	apiDao := dao.NewApiDao(os.Getenv("API_ENDPOINT"))
	playerService = service.NewPlayerService(playerDao, apiDao)

	bucketName := os.Getenv("WORDS_BUCKET")
	wordsDao := dao.NewWordsDao(bucketName)
	words, err := wordsDao.GetWords()
	if err != nil {
		fmt.Println("error loading words:", err)
		return
	}
	fmt.Println("Init: loaded", len(words), "words")
	wordsByType = words.GroupByType()
}

func handler(gameId string) error {
	// Fetch the info we need
	fmt.Println("Getting players")
	players, err := playerDao.GetPlayers(gameId)
	if err != nil {
		return fmt.Errorf("error getting players: %w\n", err)
	}
	fmt.Println("Getting game")
	game, err := gameDao.GetGame(gameId)
	if err != nil {
		return fmt.Errorf("error getting game: %w\n", err)
	}

	if players.PlayerWithHighestPoints().Points < game.TargetScore {
		// Prepare question and answer
		fmt.Println("Preparing a new question")
		wordType := model.PickRandomType()
		wordsInThisRound := wordsByType[wordType].PickRandomWords(game.OptionsPerQuestion)
		game.CorrectAnswer = wordsInThisRound.PickRandomIndex()
		game.StartTime = time.Now()
		fmt.Println("Updating game")
		err = gameDao.PutGame(game)
		if err != nil {
			return fmt.Errorf("error saving game: %w\n", err)
		}

		fmt.Println("Updating players to waiting (todo)")
		/*
		 set players.waiting
		 save players
		*/

		fmt.Println("Sending question to all players")
		err = playerService.SendQuestionToActivePlayers(players, wordsInThisRound, game.CorrectAnswer, game.SecondsPerQuestion)
		if err != nil {
			return fmt.Errorf("error sending msg to players: %w\n", err)
		}

	} else {
		// 	send game summary to all active players
		fmt.Println("todo: send game summary to all players")
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
