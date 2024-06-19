package scene

import (
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/prizelobby/ebitengine-template/core"
	"github.com/prizelobby/ebitengine-template/res"
	"github.com/prizelobby/ebitengine-template/ui"
	"github.com/prizelobby/ebitengine-template/util"
	einput "github.com/quasilyte/ebitengine-input"
)

const DRAG_SOURCE_VIEW_1 = 0
const DRAG_SOURCE_VIEW_2 = 1
const DRAG_SOURCE_DISCARD = 2

type GameUIState int

const (
	WAITING_FOR_PLAYER_MOVE GameUIState = iota
	WAITING_FOR_OPP_MOVE
	WAITING_FOR_PLAYER_ANIMIMATION
	GAME_OVER
)

type GameScene struct {
	BaseScene
	Game         *core.Game
	SelectedCard *core.Card

	Input     *einput.Handler
	PendIndex int
	PrevPend  int
	P1Score   int
	P2Score   int

	DragSprite    *ui.SphereSprite
	ViewSprite1   *ui.SphereSprite
	ViewSprite2   *ui.SphereSprite
	DiscardSprite *ui.SphereSprite
	P1Spheres     [10]*ui.SphereSprite
	P2Spheres     [10]*ui.SphereSprite
	MapSmall      *ebiten.Image
	HexMap        *ebiten.Image
	BaseTile      *ebiten.Image
	HoverTile     *ebiten.Image
	OutlineTile   *ebiten.Image

	Stroke *ui.Stroke

	HelpText string
}

const (
	PrimaryKey einput.Action = iota
	SecondaryKey
	Left
	Right
	Up
	Down
)

func NewGameScene(game *core.Game, inputHandler *einput.Handler) *GameScene {
	/*
		p1s := [10]*ui.SphereSprite{}
		for i := range 10 {
			p1s[i] = ui.NewSphereSprite(&core.Card{Value: i, Color: 0}, P2XLocs[i], P2YLocs[i])
		}*/

	return &GameScene{
		Game:        game,
		Input:       inputHandler,
		HexMap:      res.GetImage("hexmap"),
		BaseTile:    res.GetImage("hextile3"),
		HoverTile:   res.GetImage("hexoutlinegreen"),
		OutlineTile: res.GetImage("hexoutlinebroken"),
		MapSmall:    res.GetImage("circlemapsmall"),
		//P2Spheres: p1s,
		PendIndex: -1,
		HelpText:  "Click the deck to choose between 2 cards to add to your pyramid.",
	}
}

const DECK_BUTTON_X = 460
const DECK_BUTTON_Y = 260
const DECK_BUTTON_W = 120
const DECK_BUTTON_H = 146

const VIEW_ONE_X = 400
const VIEW_ONE_Y = 400
const VIEW_TWO_X = 630
const VIEW_TWO_Y = 400
const DISCARD_X = 690
const DISCARD_Y = 260

const P1StartX float64 = 800
const P1StartY float64 = 360
const P2StartX float64 = 30
const P2StartY float64 = 200

const P2YDiff = -350
const PyramidLength float64 = 93
const SQRT_3 float64 = 1.732

var P1XLocs [10]float64 = [10]float64{
	//P1StartX, P1StartX - PyramidLength/2, P1StartX + PyramidLength/2, P1StartX - PyramidLength, P1StartX, P1StartX + PyramidLength,
	//P1StartX, P1StartX - PyramidLength/2, P1StartX + PyramidLength/2, P1StartX,
	P1StartX + ui.TILE_X_OFFSET, P1StartX + ui.TILE_X_OFFSET/2, P1StartX + ui.TILE_X_OFFSET*1.5, P1StartX, P1StartX + ui.TILE_X_OFFSET, P1StartX + 2*ui.TILE_X_OFFSET,
	P1StartX + ui.TILE_X_OFFSET, P1StartX + ui.TILE_X_OFFSET/2, P1StartX + ui.TILE_X_OFFSET*1.5, P1StartX + ui.TILE_X_OFFSET,
}
var P1YLocs [10]float64 = [10]float64{
	//P1StartY, P1StartY + PyramidLength*SQRT_3/2, P1StartY + PyramidLength*SQRT_3/2, P1StartY + PyramidLength*SQRT_3, P1StartY + PyramidLength*SQRT_3, P1StartY + PyramidLength*SQRT_3,
	//P1StartY + PyramidLength*SQRT_3/3, P1StartY + PyramidLength*SQRT_3/2 + PyramidLength*SQRT_3/3, P1StartY + PyramidLength*SQRT_3/2 + PyramidLength*SQRT_3/3,
	//P1StartY + PyramidLength*4*SQRT_3/6,
	P1StartY, P1StartY + ui.TILE_Y_OFFSET, P1StartY + ui.TILE_Y_OFFSET, P1StartY + 2*ui.TILE_Y_OFFSET, P1StartY + 2*ui.TILE_Y_OFFSET, P1StartY + 2*ui.TILE_Y_OFFSET,
	P1StartY + math.Floor(ui.TILE_Y_OFFSET/2), P1StartY + math.Floor(ui.TILE_Y_OFFSET*1.5), P1StartY + math.Floor(ui.TILE_Y_OFFSET*1.5),
	P1StartY + ui.TILE_Y_OFFSET,
}

var P2XLocs [10]float64 = [10]float64{
	P2StartX + ui.TILE_X_OFFSET, P2StartX + ui.TILE_X_OFFSET*1.5, P2StartX + ui.TILE_X_OFFSET/2, P2StartX + 2*ui.TILE_X_OFFSET, P2StartX + ui.TILE_X_OFFSET, P2StartX,
	P2StartX + ui.TILE_X_OFFSET, P2StartX + ui.TILE_X_OFFSET*1.5, P2StartX + ui.TILE_X_OFFSET/2, P2StartX + ui.TILE_X_OFFSET,
}
var P2YLocs [10]float64 = [10]float64{
	P2StartY + 2*ui.TILE_Y_OFFSET, P2StartY + ui.TILE_Y_OFFSET, P2StartY + ui.TILE_Y_OFFSET, P2StartY, P2StartY, P2StartY,
	P2StartY + ui.TILE_Y_OFFSET + ui.TILE_TIP_HEIGHT, P2StartY + ui.TILE_TIP_HEIGHT, P2StartY + ui.TILE_TIP_HEIGHT,
	P2StartY + ui.TILE_TIP_HEIGHT + ui.TILE_TIP_HEIGHT,
}

var P2DrawOrder [10]int = [10]int{3, 4, 5, 1, 2, 0, 7, 8, 6, 9}

func (g *GameScene) Draw(screen *ui.ScaledScreen) {
	screen.Screen.Fill(color.RGBA{0x5b, 0xa1, 0x2d, 0xff})

	screen.DrawImage(g.MapSmall, &ebiten.DrawImageOptions{})

	//screen.DrawText(strconv.Itoa(g.PendIndex), 24, 10, 10, color.White)
	screen.DrawTextCenteredAt(g.HelpText, 24, 640, 20, color.White)
	screen.DrawTextCenteredAt("Player "+strconv.Itoa(g.Game.Turn%2+1)+"'s turn", 24, 640, 50, color.White)

	deckOpts := &ebiten.DrawImageOptions{}
	deckOpts.GeoM.Translate(DECK_BUTTON_X, DECK_BUTTON_Y)
	screen.DrawImage(g.BaseTile, deckOpts)
	screen.DrawTextCenteredAt(strconv.Itoa(len(g.Game.Deck))+"\nCards Left", 20, DECK_BUTTON_X+ui.TILE_X_OFFSET/2, DECK_BUTTON_Y+60, color.Black)

	pOpts := &ebiten.DrawImageOptions{}
	//pOpts.GeoM.Scale(1, 1)
	pOpts.GeoM.Translate(P1StartX, P1StartY)
	screen.DrawImage(g.HexMap, pOpts)
	screen.DrawTextCenteredAt("Score: "+strconv.Itoa(g.P1Score), 36, P1StartX+ui.TILE_X_OFFSET*1.5, P1StartY-30, color.White)

	pOpts2 := &ebiten.DrawImageOptions{}
	pOpts2.GeoM.Scale(1, -1)
	pOpts2.GeoM.Translate(0, 316)
	pOpts2.GeoM.Translate(P2StartX, P2StartY)
	screen.DrawImage(g.HexMap, pOpts2)
	screen.DrawTextCenteredAt("Score: "+strconv.Itoa(g.P2Score), 36, P2StartX+ui.TILE_X_OFFSET*1.5, P2StartY-30, color.White)

	screen.DrawTextCenteredAt("Discard Stack:", 30, DISCARD_X+ui.TILE_X_OFFSET/2, DISCARD_Y-50, color.White)
	screen.DrawTextCenteredAt("Deck:", 30, DECK_BUTTON_X+ui.TILE_X_OFFSET/2, DECK_BUTTON_Y-50, color.White)
	/*
		deckStartX := 400.0 - 118
		deckStartY := 100.0
		for i := range len(g.Game.Deck) {
			if i == 20 {
				deckStartX = 459 - 118
			}
			if i%5 == 0 {
				deckStartX += 118
				if i < 19 {
					deckStartY = 100
				} else {
					deckStartY = 194
				}
			} else {
				deckStartY -= 20
			}

			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(deckStartX, deckStartY)
			screen.DrawImage(g.BaseTile, opts)
		}*/

	if g.PendIndex != -1 && g.PendIndex <= 5 {
		opt := &ebiten.DrawImageOptions{}
		_, x, y := g.PyramidXYForTurn(g.PendIndex)
		opt.GeoM.Translate(x, y)
		screen.DrawImage(g.HoverTile, opt)
	}

	for i, s := range g.P1Spheres {
		if s != nil {
			s.Draw(screen)
		}
		if i == 5 {
			for j := 6; j < 9; j++ {
				if g.Game.Pyramid1.CanPlace(j) {
					opt := &ebiten.DrawImageOptions{}
					opt.GeoM.Translate(P1XLocs[j], P1YLocs[j])
					screen.DrawImage(g.OutlineTile, opt)
				}
			}
			if g.Game.Turn%2 == 0 && g.PendIndex > 5 {
				opt := &ebiten.DrawImageOptions{}
				opt.GeoM.Translate(P1XLocs[g.PendIndex], P1YLocs[g.PendIndex])
				screen.DrawImage(g.HoverTile, opt)
			}
		}
	}
	for i := range 10 {
		ii := P2DrawOrder[i]
		s := g.P2Spheres[ii]
		if s != nil {
			s.Draw(screen)
		}
		if i == 5 {
			for j := 6; j < 9; j++ {
				if g.Game.Pyramid2.CanPlace(j) {
					opt := &ebiten.DrawImageOptions{}
					opt.GeoM.Translate(P2XLocs[j], P2YLocs[j])
					screen.DrawImage(g.OutlineTile, opt)
				}
			}
			if g.Game.Turn%2 == 1 && g.PendIndex > 5 {
				opt := &ebiten.DrawImageOptions{}
				opt.GeoM.Translate(P2XLocs[g.PendIndex], P2YLocs[g.PendIndex])
				screen.DrawImage(g.HoverTile, opt)
			}
		}
	}

	if g.Game.Pyramid1.CanPlace(9) {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(P1XLocs[9], P1YLocs[9])
		screen.DrawImage(g.OutlineTile, opt)
	}
	if g.Game.Pyramid2.CanPlace(9) {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(P2XLocs[9], P2YLocs[9])
		screen.DrawImage(g.OutlineTile, opt)
	}
	if g.PendIndex == 9 {
		opt := &ebiten.DrawImageOptions{}
		_, x, y := g.PyramidXYForTurn(g.PendIndex)
		opt.GeoM.Translate(x, y)
		screen.DrawImage(g.HoverTile, opt)
	}

	if g.DiscardSprite != nil {
		g.DiscardSprite.Draw(screen)
	} else {
		screen.DrawTextCenteredAt("No cards in\ndiscard pile", 18, DISCARD_X+ui.TILE_X_OFFSET/2, DISCARD_Y+30, color.White)
	}
	if g.ViewSprite1 != nil {
		g.ViewSprite1.Draw(screen)
	}
	if g.ViewSprite2 != nil {
		g.ViewSprite2.Draw(screen)
	}
	if g.DragSprite != nil {
		g.DragSprite.Draw(screen)
	}

}

func XYinHexCell(x, y float64, Hx, Hy, Hw, Hh, Hth float64) bool {
	if !util.XYinRect(x, y, Hx, Hy, Hw, Hh) {
		return false
	}
	if util.XYinRect(x, y, Hx, Hy+Hth, Hw, Hh-2*Hth) {
		return true
	}
	if x < Hx+Hw/2 && y > -2*(x-Hx)*Hth/Hw+Hy+Hth && y < Hy+Hh-Hth+2*(x-Hx)*Hth/(Hw) {
		return true
	}
	if x >= Hx+Hw/2 && y > 2*(x-Hx-Hw/2)*Hth/Hw+Hy && y < Hy+Hh-2*(x-Hx-Hw/2)*Hth/(Hw) {
		return true
	}
	return false
}

func (g *GameScene) PyramidXYForTurn(i int) (*core.Pyramid, float64, float64) {
	if i == -1 {
		return nil, 0, 0
	}

	x, y := P1XLocs[i], P1YLocs[i]
	pyramid := g.Game.Pyramid1
	if g.Game.Turn%2 == 1 {
		x, y = P2XLocs[i], P2YLocs[i]
		pyramid = g.Game.Pyramid2
	}
	return pyramid, x, y
}

func (g *GameScene) Update() {
	/*
		if g.SelectedCard == nil {
			if g.ViewCard1 == nil {
				if g.Input.ActionIsJustPressed(PrimaryKey) {
					g.Game.DrawCards()
					g.ViewCard1 = g.Game.ViewCards[0]
					g.ViewCard2 = g.Game.ViewCards[1]
					g.ViewSprite1 = ui.NewSphereSprite(g.ViewCard1, 30, 600)
					g.ViewSprite2 = ui.NewSphereSprite(g.ViewCard2, 150, 600)
				} else if g.Input.ActionIsJustPressed(SecondaryKey) {
					if len(g.Game.Discards) > 0 {
						g.SelectedCard = g.Game.Discards[len(g.Game.Discards)-1]
						g.Game.Discards = g.Game.Discards[:len(g.Game.Discards)-1]
					}
				}
			} else {
				if g.Input.ActionIsJustPressed(Left) {
					g.SelectedCard = g.ViewCard1
					g.Game.Discards = append(g.Game.Discards, g.ViewCard2)
					g.ViewCard1 = nil
					g.ViewCard2 = nil
				} else if g.Input.ActionIsJustPressed(Right) {
					g.SelectedCard = g.ViewCard2
					g.Game.Discards = append(g.Game.Discards, g.ViewCard1)
					g.ViewCard1 = nil
					g.ViewCard2 = nil
				}
			}
		} else {
			pyramid := g.Game.Pyramid1
			if g.Game.Turn%2 == 1 {
				pyramid = g.Game.Pyramid2
			}

			if g.Input.ActionIsJustPressed(Right) {
				p := (g.PendIndex + 1) % 10
				for !pyramid.CanPlace(p) {
					p = (p + 1) % 10
				}
				g.PendIndex = p
			} else if g.Input.ActionIsJustPressed(Left) {
				p := (g.PendIndex + 9) % 10
				for !pyramid.CanPlace(p) {
					p = (p + 9) % 10
				}
				g.PendIndex = p
			} else if g.Input.ActionIsJustPressed(PrimaryKey) {
				if pyramid.CanPlace(g.PendIndex) {
					pyramid.Cards[g.PendIndex] = g.SelectedCard
					g.SelectedCard = nil
					g.Game.Turn += 1
				}
			}
		}*/

	g.PrevPend = g.PendIndex

	cx, cy := ui.AdjustedCursorPosition()
	InHex := false
	if g.DragSprite != nil {
		for i := range 10 {
			pyramid, x, y := g.PyramidXYForTurn(i)
			if XYinHexCell(cx, cy, x, y, ui.SPHERE_SIZE_X, ui.SPHERE_SIZE_Y-ui.TILE_HEIGHT, ui.TILE_TIP_HEIGHT) && pyramid.CanPlace(i) {
				g.PendIndex = i
				if g.PendIndex != g.PrevPend {
					if g.Game.Turn%2 == 0 {
						g.P1Score = pyramid.TentativeScoreWithCard(g.DragSprite.Card, g.PendIndex)
					} else {
						g.P2Score = pyramid.TentativeScoreWithCard(g.DragSprite.Card, g.PendIndex)
					}
				}

				InHex = true
				break
			}
		}
	}
	if !InHex {
		g.PendIndex = -1
		if g.PrevPend != -1 {
			g.P1Score = g.Game.Pyramid1.Score()
			g.P2Score = g.Game.Pyramid2.Score()
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.ViewSprite1 != nil && g.ViewSprite1.In(cx, cy) {
			g.Stroke = ui.NewStroke(cx, cy, g.ViewSprite1, DRAG_SOURCE_VIEW_1)
			g.DragSprite = g.ViewSprite1
			g.DragSprite.X = cx - ui.SPHERE_SIZE_X/2
			g.DragSprite.Y = cy - (ui.SPHERE_SIZE_Y-ui.TILE_HEIGHT-6)/2
		} else if g.ViewSprite2 != nil && g.ViewSprite2.In(cx, cy) {
			g.Stroke = ui.NewStroke(cx, cy, g.ViewSprite2, DRAG_SOURCE_VIEW_2)
			g.DragSprite = g.ViewSprite2
			g.DragSprite.X = cx - ui.SPHERE_SIZE_X/2
			g.DragSprite.Y = cy - (ui.SPHERE_SIZE_Y-ui.TILE_HEIGHT-6)/2
		} else if g.DiscardSprite != nil && g.DiscardSprite.In(cx, cy) {
			if g.ViewSprite1 == nil && g.ViewSprite2 == nil {
				g.Stroke = ui.NewStroke(cx, cy, g.DiscardSprite, DRAG_SOURCE_DISCARD)
				g.DragSprite = g.DiscardSprite
				g.DragSprite.X = cx - ui.SPHERE_SIZE_X/2
				g.DragSprite.Y = cy - (ui.SPHERE_SIZE_Y-ui.TILE_HEIGHT-6)/2
			} else {
				// TODO: this should hide the discard card or apply some other effect to show that it is unavailable
				g.HelpText = "You can not play the discard card after viewing from the deck this turn."
			}

		} else if util.XYinRect(cx, cy, DECK_BUTTON_X, DECK_BUTTON_Y, DECK_BUTTON_W, DECK_BUTTON_H) {
			if g.ViewSprite1 == nil && g.ViewSprite2 == nil {
				g.HelpText = "Drag one of the two cards to your pyramid."
				g.Game.DrawCards()
				g.ViewSprite1 = ui.NewSphereSprite(g.Game.ViewCards[0], VIEW_ONE_X, VIEW_ONE_Y)
				g.ViewSprite2 = ui.NewSphereSprite(g.Game.ViewCards[1], VIEW_TWO_X, VIEW_TWO_Y)
			}
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if g.Stroke != nil {
			g.Stroke.Release()
			pyramid, x, y := g.PyramidXYForTurn(g.PendIndex)
			if g.PendIndex != -1 && pyramid.CanPlace(g.PendIndex) {
				if g.Stroke.DragSourceIndex == DRAG_SOURCE_VIEW_1 {
					g.Game.PlayFromView(0, g.PendIndex)
					g.DiscardSprite = g.ViewSprite2
				} else if g.Stroke.DragSourceIndex == DRAG_SOURCE_VIEW_2 {
					g.Game.PlayFromView(1, g.PendIndex)
					g.DiscardSprite = g.ViewSprite1
				} else if g.Stroke.DragSourceIndex == DRAG_SOURCE_DISCARD {
					g.Game.PlayFromDiscard(g.PendIndex)
					if len(g.Game.Discards) > 0 {
						g.DiscardSprite = ui.NewSphereSprite(g.Game.Discards[len(g.Game.Discards)-1], DISCARD_X, DISCARD_Y)
					} else {
						g.DiscardSprite = nil
					}
				}

				g.DragSprite.X = x
				g.DragSprite.Y = y - ui.TILE_HEIGHT
				if g.DiscardSprite != nil {
					g.DiscardSprite.X = DISCARD_X
					g.DiscardSprite.Y = DISCARD_Y
					g.HelpText = "Click the deck to choose between the next 2 cards or drag the discarded card to your pyramid."
				} else {
					g.HelpText = "Click the deck to choose between 2 cards to add to your pyramid."
				}
				// this conditional is flipped because we increment the turn counter during play card
				if g.Game.Turn%2 == 0 {
					g.P2Spheres[g.PendIndex] = g.DragSprite
				} else {
					g.P1Spheres[g.PendIndex] = g.DragSprite
				}
				g.ViewSprite1 = nil
				g.ViewSprite2 = nil
				g.P1Score = g.Game.Pyramid1.Score()
				g.P2Score = g.Game.Pyramid2.Score()
			} else {
				if g.Stroke.DragSourceIndex == DRAG_SOURCE_VIEW_1 {
					g.DragSprite.X = VIEW_ONE_X
					g.DragSprite.Y = VIEW_ONE_Y
				} else if g.Stroke.DragSourceIndex == DRAG_SOURCE_VIEW_2 {
					g.DragSprite.X = VIEW_TWO_X
					g.DragSprite.Y = VIEW_TWO_Y
				} else if g.Stroke.DragSourceIndex == DRAG_SOURCE_DISCARD {
					g.DragSprite.X = DISCARD_X
					g.DragSprite.Y = DISCARD_Y
				}
			}
			g.SelectedCard = nil
			g.DragSprite = nil
		}
	}

	if g.Stroke != nil {
		g.Stroke.Update(cx, cy)
		if g.Stroke.Released {
			g.Stroke = nil
		}
	}
}
