package main

type WonderWordGame struct {
	Challenges []*Challenge
}

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

func (wg *WonderWordGame) GetRandomChallenge() *Challenge {
	return &Challenge{}
}
