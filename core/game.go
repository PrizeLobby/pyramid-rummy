package core

import (
	"math/rand"
	"strconv"
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
	State     GameState
	DrawsLeft int
}

func (g *Game) CurrentPlayer() int {
	return g.Turn % 2
}

func (g *Game) TopDiscard() *Card {
	if l := len(g.Discards); l != 0 {
		return g.Discards[l-1]
	}
	return nil
}

func (g *Game) DrawCard() *Card {
	if g.DrawsLeft == 0 {
		return nil
	}
	g.DrawsLeft -= 1
	c := g.Deck[0]
	g.Discards = append(g.Discards, c)
	g.Deck = g.Deck[1:]
	return c
}

func (g *Game) PlayCard(target int) *Card {
	c := g.Discards[len(g.Discards)-1]
	g.Discards = g.Discards[:len(g.Discards)-1]
	// UI should prevent from making illegal moves
	if g.Turn%2 == 0 {
		g.Pyramid1.Cards[target] = c
	} else {
		g.Pyramid2.Cards[target] = c
	}
	g.Turn += 1
	g.DrawsLeft = 2
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
	return c
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

var E0 [3]int = [3]int{0, 1, 3}
var E1 [3]int = [3]int{0, 2, 5}
var E2 [3]int = [3]int{3, 4, 5}
var E3 [3]int = [3]int{0, 6, 9}
var E4 [3]int = [3]int{3, 7, 9}
var E5 [3]int = [3]int{5, 8, 9}
var Edges [6][3]int = [6][3]int{E0, E1, E2, E3, E4, E5}

func (p *Pyramid) TentativeScoreWithCard(c *Card, i int) int {
	cards := p.Cards
	cards[i] = c
	score := 0
	for i := range 6 {
		d := Edges[i]
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
		d := Edges[i]
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

func (c Card) String() string {
	v := "T"
	if c.Value < 10 {
		v = strconv.Itoa(c.Value)
	}
	cs := "w"
	if c.Color == 1 {
		cs = "r"
	}
	return v + cs
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

	discards := make([]*Card, 0, 40)
	//discards = append(discards, deck[0])
	//deck = deck[1:]

	return &Game{
		Rand:      r,
		Deck:      deck,
		Discards:  discards,
		Pyramid1:  &Pyramid{Cards: [10]*Card{}},
		Pyramid2:  &Pyramid{Cards: [10]*Card{}},
		Turn:      0,
		DrawsLeft: 2,
	}
}
