package core

import (
	"math/rand"
)

type EventType int

const (
	DRAW_CARDS EventType = iota
	PLAY_CARD
)

type AgentEvent struct {
	EventType EventType
	Target    int
}

type GameAgent interface {
	GenerateMove() AgentEvent
	AcceptMove(card *Card, index int)
	RevealCard(card *Card)
}

type AgentCard struct {
	Value int
	Color int
}

type RandomAgent struct {
	PlayerNumber   int
	Rand           *rand.Rand
	ViewsRemaining int
	VisibleCard    *Card
	Pyramids       [2]*Pyramid
}

func NewAgent(playerNumber int) *RandomAgent {
	return &RandomAgent{
		Rand:           rand.New(rand.NewSource(0)),
		PlayerNumber:   playerNumber,
		Pyramids:       [2]*Pyramid{&Pyramid{}, &Pyramid{}},
		ViewsRemaining: 2,
	}
}

func (a *RandomAgent) GenerateMove() AgentEvent {
	if a.VisibleCard == nil {
		return AgentEvent{
			EventType: DRAW_CARDS,
		}
	}

	if a.ViewsRemaining > 0 {
		r := a.Rand.Intn(a.ViewsRemaining + 1)
		if r != 0 {
			return AgentEvent{
				EventType: DRAW_CARDS,
			}
		}
	}

	target := a.ChooseSlot()
	a.Pyramids[a.PlayerNumber].Cards[target] = a.VisibleCard
	a.VisibleCard = nil
	a.ViewsRemaining = 2
	return AgentEvent{
		EventType: PLAY_CARD,
		Target:    target,
	}
}

func (a *RandomAgent) ChooseSlot() int {
	seen := 0
	choice := 0
	for i := range 10 {
		if a.Pyramids[a.PlayerNumber].CanPlace(i) {
			seen += 1
			r := a.Rand.Intn(seen)
			if r == 0 {
				choice = i
			}
		}
	}
	return choice
}

func (a *RandomAgent) AcceptMove(card *Card, index int) {
	a.Pyramids[1-a.PlayerNumber].Cards[index] = card
	a.ViewsRemaining = 2
}

func (a *RandomAgent) RevealCard(card *Card) {
	a.VisibleCard = card
	a.ViewsRemaining -= 1
}
