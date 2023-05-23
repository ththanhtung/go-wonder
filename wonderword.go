package main

import (
	"math/rand"
	"strings"
	"time"
)

type Challenge struct {
	Word string
	Desc string
}

func NewChallenge(word, desc string) *Challenge {
	return &Challenge{
		Word: word,
		Desc: desc,
	}
}

var challenges = []*Challenge{
	NewChallenge("tung", "most handsome guy in the world"),
	NewChallenge("javascript", "write one bug everywhere"),
	NewChallenge("go", "google language"),
	NewChallenge("ruby", "something on the rail"),
}

type WonderWordGame struct {
	Challenge    *Challenge
	RevealedWord []string
}

func NewWonderWordGame() *WonderWordGame {
	return &WonderWordGame{}
}

func (wg *WonderWordGame) Start(){
	
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

func (wg *WonderWordGame) CheckIfWinning()bool{
	revealed := strings.Join(wg.RevealedWord, "")
	return revealed == wg.Challenge.Word
}