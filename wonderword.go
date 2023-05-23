package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

type Challenge struct {
	Word string `json:"Word"`
	Desc string `json:"Desc"`
}

func NewChallenge(word, desc string) *Challenge {
	return &Challenge{
		Word: word,
		Desc: desc,
	}
}

func LoadData() []*Challenge {
	// Read the JSON file into a byte slice
	var jsonData, err = ioutil.ReadFile("./words.json")
	if err != nil {
		panic(err)
	}

	// Unmarshal the JSON data into a slice of Challenge structs
	var challenges []*Challenge
	err = json.Unmarshal(jsonData, &challenges)
	if err != nil {
		panic(err)
	}

	return challenges
}

var challenges = LoadData()

type WonderWordGame struct {
	Challenge    *Challenge
	RevealedWord []string
}

func NewWonderWordGame() *WonderWordGame {
	return &WonderWordGame{}
}

func (wg *WonderWordGame) Start() {

	challenge := wg.GetRandomChallenge()

	revealed := make([]string, len(challenge.Word))
	for i := range revealed {
		revealed[i] = "*"
	}

	wg.Challenge = challenge
	wg.RevealedWord = revealed
}

func (wg *WonderWordGame) GetRandomChallenge() *Challenge {
	rand.Seed(time.Now().UnixNano())
	randWordIndex := rand.Intn(len(challenges) - 1)
	return challenges[randWordIndex]
}

func (wg *WonderWordGame) ScoreCalculator(guessChar string) (int, string) {
	guess := strings.ToLower(guessChar)
	for i, c := range wg.Challenge.Word {
		if guess == string(c) {
			wg.RevealedWord[i] = guess
		}
	}
	score := 100 * strings.Count(wg.Challenge.Word, guess)
	return score, strings.Join(wg.RevealedWord, "")
}

func (wg *WonderWordGame) CheckIfEndGame() bool {
	revealed := strings.Join(wg.RevealedWord, "")
	return revealed == wg.Challenge.Word
}
