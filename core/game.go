package core

import (
	"math/rand"
	"time"
)

type GameState int

const (
	IN_PROGRESS GameState = iota
	P1_WIN
	P2_WIN
	DRAW
)

type SourceLocation int

const (
	VIEW_ONE SourceLocation = iota
	VIEW_TWO
	DISCARD
)

type Game struct {
	Rand      *rand.Rand
	Deck      []*Card
	Discards  []*Card
	Pyramid1  *Pyramid
	Pyramid2  *Pyramid
	Turn      int
	ViewCards [2]*Card
	State     GameState
}

func (g *Game) DrawCards() {
	g.ViewCards[0] = g.Deck[0]
	g.ViewCards[1] = g.Deck[1]
	g.Deck = g.Deck[2:]
}

// UI should prevent taking discard if player chooses to draw
func (g *Game) TakeDiscard() *Card {
	card := g.Discards[len(g.Discards)-1]
	g.Discards = g.Discards[:len(g.Discards)-1]
	return card
}

// index must be 0 or 1
// ui should prevent other values
func (g *Game) TakeFromView(index int) *Card {
	c := g.ViewCards[index]
	g.Discards = append(g.Discards, g.ViewCards[1-index])
	g.ViewCards[0] = nil
	g.ViewCards[1] = nil
	return c
}

func (g *Game) PlayFromView(index int, target int) {
	g.PlayCard(g.TakeFromView(index), target)
}

func (g *Game) PlayFromDiscard(target int) {
	g.PlayCard(g.TakeDiscard(), target)
}

func (g *Game) PlayCard(c *Card, target int) {
	// UI should prevent from making illegal moves
	if g.Turn%2 == 0 {
		g.Pyramid1.Cards[target] = c
	} else {
		g.Pyramid2.Cards[target] = c
		//g.Pyramid1.Cards[target] = c
	}
	g.Turn += 1
	if g.Turn == 20 {
		s1 := g.Pyramid1.Score()
		s2 := g.Pyramid2.Score()
		if s1 > s2 {
			g.State = P1_WIN
		} else if s2 > s1 {
			g.State = P2_WIN
		} else {
			g.State = DRAW
		}
	}
}

/*
	      0

		  6

	 1	        2
	      9
	  7       8

3         4          5
*/
type Pyramid struct {
	Cards [10]*Card
}

var D0 [3]int = [3]int{0, 1, 3}
var D1 [3]int = [3]int{0, 2, 5}
var D2 [3]int = [3]int{3, 4, 5}
var D3 [3]int = [3]int{0, 6, 9}
var D4 [3]int = [3]int{3, 7, 9}
var D5 [3]int = [3]int{5, 8, 9}
var Diagonals [6][3]int = [6][3]int{D0, D1, D2, D3, D4, D5}

func (p *Pyramid) TentativeScoreWithCard(c *Card, i int) int {
	cards := p.Cards
	cards[i] = c
	score := 0
	for i := range 6 {
		d := Diagonals[i]
		if cards[d[0]] == nil || cards[d[1]] == nil || cards[d[2]] == nil {
			continue
		}
		for j := range 3 {
			if cards[d[j]].Color != cards[d[(j+1)%3]].Color && cards[d[j]].Color != cards[d[(j+2)%3]].Color {
				score += cards[d[j]].Value
			}
		}
	}
	return score
}

func (p *Pyramid) Score() int {
	score := 0
	for i := range 6 {
		d := Diagonals[i]
		if p.Cards[d[0]] == nil || p.Cards[d[1]] == nil || p.Cards[d[2]] == nil {
			continue
		}
		for j := range 3 {
			if p.Cards[d[j]].Color != p.Cards[d[(j+1)%3]].Color && p.Cards[d[j]].Color != p.Cards[d[(j+2)%3]].Color {
				score += p.Cards[d[j]].Value
			}
		}
	}
	return score
}

func (p *Pyramid) CanPlace(i int) bool {
	if p.Cards[i] != nil {
		return false
	}

	if i == 6 {
		return p.Cards[0] != nil && p.Cards[1] != nil && p.Cards[2] != nil
	} else if i == 7 {
		return p.Cards[1] != nil && p.Cards[3] != nil && p.Cards[4] != nil
	} else if i == 8 {
		return p.Cards[2] != nil && p.Cards[4] != nil && p.Cards[5] != nil
	} else if i == 9 {
		return p.Cards[6] != nil && p.Cards[7] != nil && p.Cards[8] != nil
	}
	return true
}

type Card struct {
	Value int
	Color int
}

func NewGame() *Game {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	deck := make([]*Card, 40)
	for i := range 10 {
		deck[4*i] = &Card{Value: i + 1, Color: 0}
		deck[4*i+1] = &Card{Value: i + 1, Color: 0}
		deck[4*i+2] = &Card{Value: i + 1, Color: 1}
		deck[4*i+3] = &Card{Value: i + 1, Color: 1}
	}
	r.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })

	return &Game{
		Rand:     r,
		Deck:     deck,
		Discards: make([]*Card, 0, 40),
		Pyramid1: &Pyramid{Cards: [10]*Card{}},
		Pyramid2: &Pyramid{Cards: [10]*Card{}},
		Turn:     0,
	}
}
