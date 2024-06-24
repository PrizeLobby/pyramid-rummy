package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/prizelobby/pyramid-rummy/res"
)

type RulesComponent struct {
	EdgeGuide *ebiten.Image
}

func NewRulesComponent() *RulesComponent {
	return &RulesComponent{
		EdgeGuide: res.GetImage("edgeguide"),
	}
}

func (r *RulesComponent) Update() {

}

func (r *RulesComponent) Draw(screen *ScaledScreen) {
	screen.DrawText(
		`Pyramid rummy is played with a 40 card deck. The cards are split into two suits, purple leaves and yellow nibs. 
Each suit contains 2 copies of the values 1 to 10.

Play cards to build the highest scoring pyramid. You can choose to play the revealed card on your pyramid or draw a 
new card. You may draw up to 2 cards each turn, at which you will be forced to play the most recently revealed card.

Scoring is based on the six edges of the pyramid. Each edge consists of three cards. If all three cards are the same
color, the score for that edge is 0. Otherwise, the score is equal to the value of the card that is a different color
than the other two. Your total score is the sum of the scores of the six edges.
`,
		24, 15, 15, color.White)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(1, 1)
	opts.GeoM.Translate(640-(832/2), 310)
	screen.DrawImage(r.EdgeGuide, opts)

	screen.DrawTextCenteredAt("Click anywhere to return", 18, 640, 690, color.White)
}
