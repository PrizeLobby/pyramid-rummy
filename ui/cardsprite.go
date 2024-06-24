package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/prizelobby/pyramid-rummy/core"
	"github.com/prizelobby/pyramid-rummy/res"
	"github.com/prizelobby/pyramid-rummy/util"
)

const TILE_SIZE_X = 120
const TILE_SIZE_Y = 170
const TILE_HEIGHT = 20
const TILE_X_OFFSET = 118
const TILE_Y_OFFSET = 95
const TILE_TIP_HEIGHT = 30

const DISPLAY_TYPE_REGULAR = 0
const DISPLAY_TYPE_BOTTOM = 1
const DISPLAY_TYPE_LEFT = 2
const DISPLAY_TYPE_RIGHT = 3

type CardSprite struct {
	Card        *core.Card
	Tiles       *Tileset
	Shadow      *ebiten.Image
	Shadow2     *ebiten.Image
	X, Y        float64
	ShadowType  int //0 == regular shadow, 1 == no shadow, 2 == shadow on stack
	DisplayType int //0 == top, 1 == bottom, 2 == left, 3 == right
}

func NewCardSprite(c *core.Card, x, y float64) *CardSprite {
	return &CardSprite{
		Card: c, X: x, Y: y,
		Tiles:      NewTileset(res.GetImage("hextiletileset"), 120, 146),
		Shadow:     res.GetImage("shadow"),
		Shadow2:    res.GetImage("shadow2"),
		ShadowType: 0,
	}
}

func (c *CardSprite) Draw(screen *ScaledScreen) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(c.X, c.Y)
	screen.DrawImage(c.Tiles.TileAtIJ(c.Card.Value-1, c.Card.Color*4+c.DisplayType), opts)

	if c.ShadowType == 0 {
		opts.GeoM.Reset()
		opts.GeoM.Translate(c.X, c.Y)
		screen.DrawImage(c.Shadow, opts)
	} else if c.ShadowType == 2 {
		opts.GeoM.Reset()
		opts.GeoM.Translate(c.X, c.Y)
		opts.ColorScale.ScaleAlpha(0.8)
		screen.DrawImage(c.Shadow2, opts)
	}
}

func (c *CardSprite) In(x, y float64) bool {
	return util.XYinRect(x, y, c.X, c.Y, TILE_SIZE_X, TILE_SIZE_Y)
}

func (c *CardSprite) MoveBy(dx, dy float64) {
	c.X += dx
	c.Y += dy
}

func (c *CardSprite) SetPos(x, y float64) {
	c.X = x
	c.Y = y
}
