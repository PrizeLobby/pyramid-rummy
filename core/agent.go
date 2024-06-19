package core

import (
	"math/rand"
)

type EventType int

const (
	DRAW_CARDS EventType = iota
	PLAY_DISCARD
	PLAY_VIEW1
	PLAY_VIEW2
)

type AgentEvent struct {
	EventType EventType
	Target    int
}

type GameAgent interface {
	GenerateMove() AgentEvent
}

type AgentCard struct {
	Value int
	Color int
}

type Agent struct {
	PlayerNumber int
	Rand         *rand.Rand
}

func NewAgent() *Agent {
	return &Agent{
		Rand: rand.New(rand.NewSource(0)),
	}
}

func (a *Agent) GenerateMove() AgentEvent {
	return AgentEvent{}
}

func (a *Agent) AcceptMove(card *Card, index int) {

}
