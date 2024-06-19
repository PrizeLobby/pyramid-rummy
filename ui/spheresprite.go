package ui

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/prizelobby/ebitengine-template/core"
	"github.com/prizelobby/ebitengine-template/res"
	"github.com/prizelobby/ebitengine-template/util"
)

const SPHERE_SIZE_X = 120
const SPHERE_SIZE_Y = 170
const TILE_HEIGHT = 20
const TILE_X_OFFSET = 118
const TILE_Y_OFFSET = 95
const TILE_TIP_HEIGHT = 30

type SphereSprite struct {
	Card *core.Card
	Img  *ebiten.Image
	X, Y float64
}

func NewSphereSprite(c *core.Card, x, y float64) *SphereSprite {
	var img *ebiten.Image
	if c.Color == 0 {
		img = res.GetImage("hextile")
	} else {
		img = res.GetImage("hextile2")
	}
	return &SphereSprite{
		Card: c, X: x, Y: y,
		Img: img,
	}
}

func (s *SphereSprite) Update() {

}

func (s *SphereSprite) Draw(screen *ScaledScreen) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(s.X, s.Y)
	screen.DrawImage(s.Img, opts)
	var col color.Color = color.White
	if s.Card.Color == 1 {
		col = color.RGBA{0x78, 0x21, 0x4b, 0xFF}
	}
	screen.DrawTextCenteredAt(strconv.Itoa(s.Card.Value), 32, s.X+SPHERE_SIZE_X/2, s.Y+SPHERE_SIZE_Y/2-10, col)
}

func (s *SphereSprite) In(x, y float64) bool {
	return util.XYinRect(x, y, s.X, s.Y, SPHERE_SIZE_X, SPHERE_SIZE_Y)
}

func (s *SphereSprite) MoveBy(dx, dy float64) {
	s.X += dx
	s.Y += dy
}
