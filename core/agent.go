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
	SetVisibleCard(card *Card)
}

func CardToIndex(c *Card) int {
	return c.Color*20 + c.Copy*10 + c.Value - 1
}

func IndexToCard(i int) *Card {
	return &Card{
		Value: i%10 + 1,
		Color: i / 20,
	}
}

func NewShuffledDeck(r *rand.Rand) [40]*Card {
	d := NewDeck()
	r.Shuffle(40, func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
	return d
}

type RandomAgent struct {
	PlayerNumber   int
	Rand           *rand.Rand
	ViewsRemaining int
	VisibleCard    *Card
	Pyramids       [2]*Pyramid
}

func NewRandomAgent(playerNumber int) *RandomAgent {
	return &RandomAgent{
		Rand:           rand.New(rand.NewSource(0)),
		PlayerNumber:   playerNumber,
		Pyramids:       [2]*Pyramid{&Pyramid{}, &Pyramid{}},
		ViewsRemaining: 2,
	}
}

func (a *RandomAgent) SetVisibleCard(c *Card) {
	a.VisibleCard = c
}

func (a *RandomAgent) GenerateMove() AgentEvent {
	if a.VisibleCard == nil {
		a.ViewsRemaining -= 1
		//fmt.Printf("Agent %d: Drawing card\n", a.PlayerNumber+1)
		return AgentEvent{
			EventType: DRAW_CARDS,
		}
	}

	if a.ViewsRemaining > 0 {
		r := a.Rand.Intn(a.ViewsRemaining + 1)
		if r != 0 {
			a.ViewsRemaining -= 1
			//fmt.Printf("Agent %d: Drawing card\n", a.PlayerNumber+1)
			return AgentEvent{
				EventType: DRAW_CARDS,
			}
		}
	}

	target := a.ChooseSlot()
	a.Pyramids[a.PlayerNumber].Cards[target] = a.VisibleCard
	//fmt.Printf("Agent %d: Playing card %s at %d\n", a.PlayerNumber+1, a.VisibleCard.String(), target)
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

type SampleAgent struct {
	Orientation    int
	PlayerNumber   int
	Rand           *rand.Rand
	DrawsRemaining int
	CardsPlayed    int
	VisibleCard    *Card
	Pyramids       [2]*Pyramid
	SeenCards      [40]bool

	Strategy int
}

func NewSampleAgent(playerNumber int) *SampleAgent {
	r := rand.New(rand.NewSource(0))
	orientation := rand.Intn(6)
	return &SampleAgent{
		Rand:           r,
		Orientation:    orientation,
		PlayerNumber:   playerNumber,
		Pyramids:       [2]*Pyramid{&Pyramid{}, &Pyramid{}},
		DrawsRemaining: 2,
	}
}

func (a *SampleAgent) RecordMove(target int) AgentEvent {
	//fmt.Printf("Agent %d: Playing card "+a.VisibleCard.String()+" at target %d\n", a.PlayerNumber+1, target)
	a.Pyramids[a.PlayerNumber].Cards[target] = a.VisibleCard
	a.VisibleCard = nil
	a.DrawsRemaining = 2
	a.CardsPlayed += 1
	return AgentEvent{
		EventType: PLAY_CARD,
		Target:    target,
	}
}

func (a *SampleAgent) RecordDraw() AgentEvent {
	a.DrawsRemaining -= 1
	return AgentEvent{
		EventType: DRAW_CARDS,
	}
}

func (a *SampleAgent) AvailableSlots() []int {
	if a.CardsPlayed == 0 {
		switch a.Orientation {
		case 0:
			return []int{0, 1}
		case 1:
			return []int{0, 2}
		case 2:
			return []int{3, 1}
		case 3:
			return []int{3, 4}
		case 4:
			return []int{5, 2}
		case 5:
		default:
			return []int{5, 4}
		}
	} else if a.CardsPlayed == 1 {
		switch a.Orientation {
		case 0:
			if !a.Pyramids[a.PlayerNumber].CanPlace(0) {
				return []int{1, 3, 4}
			} else {
				return []int{0, 3, 4}
			}
		case 1:
			if !a.Pyramids[a.PlayerNumber].CanPlace(0) {
				return []int{2, 4, 5}
			} else {
				return []int{0, 4, 5}
			}
		case 2:
			if !a.Pyramids[a.PlayerNumber].CanPlace(3) {
				return []int{0, 1, 2}
			} else {
				return []int{3, 0, 2}
			}
		case 3:
			if !a.Pyramids[a.PlayerNumber].CanPlace(3) {
				return []int{4, 2, 5}
			} else {
				return []int{3, 2, 5}
			}
		case 4:
			if !a.Pyramids[a.PlayerNumber].CanPlace(5) {
				return []int{2, 0, 1}
			} else {
				return []int{0, 1, 5}
			}
		case 5:
		default:
			if !a.Pyramids[a.PlayerNumber].CanPlace(5) {
				return []int{1, 3, 4}
			} else {
				return []int{5, 1, 3}
			}
		}
	} else if a.CardsPlayed == 9 {
		return []int{9}
	}
	available := make([]int, 0, 10)
	for i := range 10 {
		if a.Pyramids[a.PlayerNumber].CanPlace(i) {
			available = append(available, i)
		}
	}
	return available
}

func (a *SampleAgent) GenerateMove() AgentEvent {
	return a.GenerateMoveB(100, 20)
}

func (a *SampleAgent) RandomUnseenIndexAndCard() (int, *Card) {
	randIndex := a.Rand.Intn(40)
	for a.SeenCards[randIndex] {
		randIndex = (randIndex + 1) % 40
	}
	return randIndex, IndexToCard(randIndex)
}

func (a *SampleAgent) GenerateMoveB(iterations, drawIterations int) AgentEvent {
	if a.VisibleCard == nil {
		return a.RecordDraw()
	}
	p := a.Pyramids[a.PlayerNumber]
	slots := a.AvailableSlots()

	emptySlots := []int{}
	for i := range 10 {
		if p.Cards[i] == nil {
			emptySlots = append(emptySlots, i)
		}
	}
	samples := make([][10]*Card, iterations)
	for i := range iterations {
		//fmt.Printf("iteration %d\n", i)
		samples[i] = p.Cards
		for _, es := range emptySlots {
			var randIndex int
			randIndex, samples[i][es] = a.RandomUnseenIndexAndCard()
			a.SeenCards[randIndex] = true
		}
		for _, es := range emptySlots {
			a.SeenCards[CardToIndex(samples[i][es])] = false
		}
		//fmt.Println(samples[i])
	}

	slotScores := make([][]int, iterations) //[iterations][]int{}
	for i := range iterations {
		slotScores[i] = make([]int, len(slots))
	}
	tempPyrmaid := &Pyramid{}
	for sI, slot := range slots {
		for i := range iterations {
			tempPyrmaid.Cards = samples[i]
			score := tempPyrmaid.TentativeScoreWithCard(a.VisibleCard, slot)
			slotScores[i][sI] = score

		}
	}
	slotIndexResults := make([]int, len(slots))
	for i := range iterations {
		bestSlotIteration := 0
		bestScore := -1
		for j := range len(slots) {
			if slotScores[i][j] > bestScore {
				bestSlotIteration = j
				bestScore = slotScores[i][j]
			}
		}
		slotIndexResults[bestSlotIteration] += 1
	}
	bestSlotIndex := 0
	bestResults := -1
	for i, s := range slotIndexResults {
		if s > bestResults {
			bestSlotIndex = i
			bestResults = s
		}
	}

	if a.DrawsRemaining == 0 {
		return a.RecordMove(slots[bestSlotIndex])
	}

	pendScores := make([]int, iterations)
	for i := range iterations {
		iterScore := 0
		tempPyrmaid.Cards = samples[i]
		for range drawIterations {
			slot := slots[a.Rand.Intn(len(slots))]
			// this could generate a duplicate thats already in the pyramid
			// but it should be good enough
			_, randCard := a.RandomUnseenIndexAndCard()
			iterScore += tempPyrmaid.TentativeScoreWithCard(randCard, slot)
		}
		pendScores[i] = iterScore
	}
	drawBetter := 0
	for i := range iterations {
		drawScore := pendScores[i]
		visScore := slotScores[i][bestSlotIndex]
		if drawScore > visScore*drawIterations {
			drawBetter += 1
		}
	}
	if drawBetter > iterations/2 {
		return a.RecordDraw()
	}
	return a.RecordMove(slots[bestSlotIndex])
}

func (a *SampleAgent) AcceptMove(card *Card, index int) {
	a.Pyramids[1-a.PlayerNumber].Cards[index] = card
	a.DrawsRemaining = 2
}

func (a *SampleAgent) RevealCard(card *Card) {
	a.VisibleCard = card
	a.SeenCards[CardToIndex(card)] = true
}

func (a *SampleAgent) SetVisibleCard(c *Card) {
	a.VisibleCard = c
}
