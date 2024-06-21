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

type GameUIState int

const (
	WAITING_FOR_PLAYER_MOVE GameUIState = iota
	WAITING_FOR_OPP_MOVE
	WAITING_FOR_PLAYER_ANIMIMATION
	GAME_OVER
)

type GameScene struct {
	BaseScene
	UIState GameUIState

	Game         *core.Game
	SelectedCard *core.Card

	Agents   [2]core.GameAgent
	moveChan chan core.AgentEvent

	Input     *einput.Handler
	PendIndex int
	PrevPend  int
	P1Score   int
	P2Score   int

	DragSprite    *ui.CardSprite
	DiscardSprite *ui.CardSprite
	P1Spheres     [10]*ui.CardSprite
	P2Spheres     [10]*ui.CardSprite
	MapSmall      *ebiten.Image
	HexMap        *ebiten.Image
	BaseTile      *ebiten.Image
	HoverTile     *ebiten.Image
	OutlineTile   *ebiten.Image
	Shadow        *ebiten.Image

	Stroke *ui.Stroke

	ActiveAnimation ui.Anim
	AnimationQueue  []ui.Anim

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
	return &GameScene{
		Game:          game,
		Input:         inputHandler,
		DiscardSprite: ui.NewCardSprite(game.TopDiscard(), DISCARD_X, DISCARD_Y),
		HexMap:        res.GetImage("hexmap"),
		BaseTile:      res.GetImage("basetile"),
		HoverTile:     res.GetImage("hexoutlinegreen"),
		OutlineTile:   res.GetImage("hexoutlinebroken"),
		MapSmall:      res.GetImage("circlemapsmall"),
		Shadow:        res.GetImage("shadow"),
		moveChan:      make(chan core.AgentEvent, 1),
		PendIndex:     -1,
		HelpText:      "Drag the open card to your pyramid or click the deck to reveal a new card.",
		Agents:        [2]core.GameAgent{nil, core.NewAgent(1)},
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
const P2StartY float64 = 360

var P1XLocs [10]float64 = [10]float64{
	P1StartX + ui.TILE_X_OFFSET, P1StartX + ui.TILE_X_OFFSET/2, P1StartX + ui.TILE_X_OFFSET*1.5, P1StartX, P1StartX + ui.TILE_X_OFFSET, P1StartX + 2*ui.TILE_X_OFFSET,
	P1StartX + ui.TILE_X_OFFSET, P1StartX + ui.TILE_X_OFFSET/2, P1StartX + ui.TILE_X_OFFSET*1.5, P1StartX + ui.TILE_X_OFFSET,
}
var P1YLocs [10]float64 = [10]float64{
	P1StartY, P1StartY + ui.TILE_Y_OFFSET, P1StartY + ui.TILE_Y_OFFSET, P1StartY + 2*ui.TILE_Y_OFFSET, P1StartY + 2*ui.TILE_Y_OFFSET, P1StartY + 2*ui.TILE_Y_OFFSET,
	P1StartY + math.Floor(ui.TILE_Y_OFFSET/2), P1StartY + math.Floor(ui.TILE_Y_OFFSET*1.5), P1StartY + math.Floor(ui.TILE_Y_OFFSET*1.5),
	P1StartY + ui.TILE_Y_OFFSET - 1,
}

var P2XLocs [10]float64 = [10]float64{
	P2StartX + ui.TILE_X_OFFSET, P2StartX + ui.TILE_X_OFFSET/2, P2StartX + ui.TILE_X_OFFSET*1.5, P2StartX, P2StartX + ui.TILE_X_OFFSET, P2StartX + 2*ui.TILE_X_OFFSET,
	P2StartX + ui.TILE_X_OFFSET, P2StartX + ui.TILE_X_OFFSET/2, P2StartX + ui.TILE_X_OFFSET*1.5, P2StartX + ui.TILE_X_OFFSET,
}
var P2YLocs [10]float64 = [10]float64{
	P1StartY, P1StartY + ui.TILE_Y_OFFSET, P1StartY + ui.TILE_Y_OFFSET, P1StartY + 2*ui.TILE_Y_OFFSET, P1StartY + 2*ui.TILE_Y_OFFSET, P1StartY + 2*ui.TILE_Y_OFFSET,
	P1StartY + math.Floor(ui.TILE_Y_OFFSET/2), P1StartY + math.Floor(ui.TILE_Y_OFFSET*1.5), P1StartY + math.Floor(ui.TILE_Y_OFFSET*1.5),
	P1StartY + ui.TILE_Y_OFFSET - 1,
}

func (g *GameScene) Draw(screen *ui.ScaledScreen) {
	screen.Screen.Fill(color.RGBA{0x5b, 0xa1, 0x2d, 0xff})

	screen.DrawImage(g.MapSmall, &ebiten.DrawImageOptions{})

	//screen.DrawText(strconv.Itoa(g.PendIndex), 24, 10, 10, color.White)
	screen.DrawTextCenteredAt(g.HelpText, 24, 640, 20, color.White)
	if g.UIState != GAME_OVER {
		screen.DrawTextCenteredAt("Player "+strconv.Itoa(g.Game.Turn%2+1)+"'s turn", 24, 640, 50, color.White)
	} else {
		if g.Game.State == core.P1_WIN {
			screen.DrawTextCenteredAt("P1 Wins", 24, 640, 50, color.White)
		} else if g.Game.State == core.P2_WIN {
			screen.DrawTextCenteredAt("P2 Wins", 24, 640, 50, color.White)
		} else if g.Game.State == core.DRAW {
			screen.DrawTextCenteredAt("Draw", 24, 640, 50, color.White)
		}
	}

	deckOpts := &ebiten.DrawImageOptions{}
	deckOpts.GeoM.Translate(DECK_BUTTON_X, DECK_BUTTON_Y)
	screen.DrawImage(g.BaseTile, deckOpts)
	deckShadowOpts := &ebiten.DrawImageOptions{}
	deckShadowOpts.GeoM.Translate(DECK_BUTTON_X, DECK_BUTTON_Y)
	screen.DrawImage(g.Shadow, deckShadowOpts)
	//screen.DrawTextCenteredAt(strconv.Itoa(len(g.Game.Deck))+"\nCards Left", 20, DECK_BUTTON_X+ui.TILE_X_OFFSET/2, DECK_BUTTON_Y+60, color.Black)

	dOpts := &ebiten.DrawImageOptions{}
	dOpts.GeoM.Translate(DISCARD_X, DISCARD_Y+ui.TILE_HEIGHT)
	screen.DrawImage(g.OutlineTile, dOpts)

	if len(g.Game.Discards) < 2 {
		screen.DrawTextCenteredAt("No prior\ncards in stack", 18, DISCARD_X+ui.TILE_X_OFFSET/2, DISCARD_Y+70, color.White)
	} else {
		screen.DrawTextCenteredAt("Prior card:\n"+g.Game.Discards[len(g.Game.Discards)-2].String(), 18, DISCARD_X+ui.TILE_X_OFFSET/2, DISCARD_Y+70, color.White)
	}
	if g.DiscardSprite != nil {
		g.DiscardSprite.Draw(screen)
	}

	pOpts := &ebiten.DrawImageOptions{}
	pOpts.GeoM.Translate(P1StartX, P1StartY)
	screen.DrawImage(g.HexMap, pOpts)
	screen.DrawTextCenteredAt("Score: "+strconv.Itoa(g.P1Score), 36, P1StartX+ui.TILE_X_OFFSET*1.5, P1StartY-40, color.White)

	pOpts2 := &ebiten.DrawImageOptions{}
	pOpts2.GeoM.Translate(P2StartX, P1StartY)
	screen.DrawImage(g.HexMap, pOpts2)
	screen.DrawTextCenteredAt("Score: "+strconv.Itoa(g.P2Score), 36, P2StartX+ui.TILE_X_OFFSET*1.5, P2StartY-40, color.White)

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
	for i, s := range g.P2Spheres {
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

	select {
	case m := <-g.moveChan:
		if m.EventType == core.DRAW_CARDS {
			c := g.Game.DrawCard()
			g.Agents[g.Game.Turn%2].RevealCard(c)

			// TODO: this needs to be done after the 30 tick wait
			g.DiscardSprite = ui.NewCardSprite(c, DECK_BUTTON_X, DECK_BUTTON_Y)
			g.AnimationQueue = append(g.AnimationQueue, ui.NewBlockingAnim(30), ui.NewLinearPathAnimator(g.DiscardSprite, 35,
				ui.Location{X: DECK_BUTTON_X, Y: DECK_BUTTON_Y},
				ui.Location{X: DISCARD_X, Y: DISCARD_Y}, ui.EaseOutCubic, func() {
					go func() { g.moveChan <- g.Agents[g.Game.Turn%2].GenerateMove() }()
				}))

		} else if m.EventType == core.PLAY_CARD {
			g.Game.PlayCard(m.Target)
			complete := func() {
				nextCard := g.Game.TopDiscard()
				if nextCard != nil {
					g.DiscardSprite = ui.NewCardSprite(nextCard, DISCARD_X, DISCARD_Y)
				}
				g.P1Score = g.Game.Pyramid1.Score()
				g.P2Score = g.Game.Pyramid2.Score()
				if g.Game.State == core.IN_PROGRESS {
					if g.Agents[g.Game.Turn%2] == nil {
						g.UIState = WAITING_FOR_PLAYER_MOVE
					} else {
						go func() { g.moveChan <- g.Agents[g.Game.Turn%2].GenerateMove() }()
					}
				} else {
					g.UIState = GAME_OVER
				}
			}

			if g.Game.Turn%2 == 0 {
				g.P2Spheres[m.Target] = g.DiscardSprite // ui.NewCardSprite(card, P2XLocs[m.Target], P2YLocs[m.Target])
				g.DiscardSprite = nil
				g.AnimationQueue = append(g.AnimationQueue, ui.NewBlockingAnim(30), ui.NewLinearPathAnimator(g.P2Spheres[m.Target], 50,
					ui.Location{X: DISCARD_X, Y: DISCARD_Y},
					ui.Location{X: P2XLocs[m.Target], Y: P2YLocs[m.Target] - ui.TILE_HEIGHT}, ui.EaseOutCubic, complete))
			} else {
				g.P1Spheres[m.Target] = g.DiscardSprite // ui.NewCardSprite(card, P2XLocs[m.Target], P2YLocs[m.Target])
				g.DiscardSprite = nil
				g.AnimationQueue = append(g.AnimationQueue, ui.NewBlockingAnim(30), ui.NewLinearPathAnimator(g.P1Spheres[m.Target], 50,
					ui.Location{X: DISCARD_X, Y: DISCARD_Y},
					ui.Location{X: P1XLocs[m.Target], Y: P1YLocs[m.Target] - ui.TILE_HEIGHT}, ui.EaseOutCubic, complete))
			}
		}
	default:
	}

	if g.ActiveAnimation != nil {
		if g.ActiveAnimation.IsFinished() {
			g.ActiveAnimation = nil
			return
		} else {
			g.ActiveAnimation.Update()
		}
	} else {
		if len(g.AnimationQueue) > 0 {
			g.ActiveAnimation = g.AnimationQueue[0]
			g.AnimationQueue = g.AnimationQueue[1:]
		} else {
			//b.isAnimating = false
			//b.CompleteMove(b.animatingMove)
		}
	}

	if g.UIState == WAITING_FOR_PLAYER_MOVE {
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
			if g.DiscardSprite != nil && g.DiscardSprite.In(cx, cy) {
				g.Stroke = ui.NewStroke(cx, cy, g.DiscardSprite, 0)
				g.DragSprite = g.DiscardSprite
				g.DragSprite.ShadowType = 1
				g.DragSprite.X = cx - ui.SPHERE_SIZE_X/2
				g.DragSprite.Y = cy - (ui.SPHERE_SIZE_Y-ui.TILE_HEIGHT-6)/2
			} else if util.XYinRect(cx, cy, DECK_BUTTON_X, DECK_BUTTON_Y, DECK_BUTTON_W, DECK_BUTTON_H) {
				if g.Game.DrawsLeft == 0 {
					g.HelpText = "You have 0 draws remaining this turn."
				} else {
					c := g.Game.DrawCard()
					g.DiscardSprite = ui.NewCardSprite(c, DECK_BUTTON_X, DECK_BUTTON_Y)
					g.AnimationQueue = append(g.AnimationQueue, ui.NewLinearPathAnimator(g.DiscardSprite, 25,
						ui.Location{X: DECK_BUTTON_X, Y: DECK_BUTTON_Y},
						ui.Location{X: DISCARD_X, Y: DISCARD_Y}, ui.EaseOutCubic, nil))
				}
			}
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			if g.Stroke != nil {
				g.Stroke.Release()
				pyramid, x, y := g.PyramidXYForTurn(g.PendIndex)
				if g.PendIndex != -1 && pyramid.CanPlace(g.PendIndex) {
					g.Game.PlayCard(g.PendIndex)
					if len(g.Game.Discards) > 0 {
						g.DiscardSprite = ui.NewCardSprite(g.Game.TopDiscard(), DISCARD_X, DISCARD_Y)
						g.DiscardSprite.X = DISCARD_X
						g.DiscardSprite.Y = DISCARD_Y
						g.HelpText = "Drag the open card to your pyramid or click the deck to reveal a new card."
					} else {
						g.DiscardSprite = nil
						g.HelpText = "Click the deck to choose reveal a card."
					}
					g.DragSprite.X = x
					g.DragSprite.Y = y - ui.TILE_HEIGHT
					if g.PendIndex < 5 {
						g.DragSprite.ShadowType = 0
					} else {
						g.DragSprite.ShadowType = 2
					}
					// this conditional is flipped because we increment the turn counter during play card
					var sprites [10]*ui.CardSprite
					if g.Game.Turn%2 == 0 {
						g.P2Spheres[g.PendIndex] = g.DragSprite
						sprites = g.P2Spheres
					} else {
						g.P1Spheres[g.PendIndex] = g.DragSprite
						sprites = g.P1Spheres
					}
					if g.PendIndex == 6 {
						sprites[1].DisplayType = ui.DISPLAY_TYPE_LEFT
						sprites[2].DisplayType = ui.DISPLAY_TYPE_RIGHT
					} else if g.PendIndex == 7 {
						if sprites[6] != nil {
							sprites[1].DisplayType = ui.DISPLAY_TYPE_LEFT
						}
						sprites[3].DisplayType = ui.DISPLAY_TYPE_LEFT
						if sprites[8] != nil {
							sprites[4].DisplayType = ui.DISPLAY_TYPE_BOTTOM
						} else {
							sprites[4].DisplayType = ui.DISPLAY_TYPE_RIGHT
						}
					} else if g.PendIndex == 8 {
						if sprites[7] != nil {
							sprites[2].DisplayType = ui.DISPLAY_TYPE_RIGHT
						}
						if sprites[7] != nil {
							sprites[4].DisplayType = ui.DISPLAY_TYPE_BOTTOM
						} else {
							sprites[4].DisplayType = ui.DISPLAY_TYPE_LEFT
						}
						sprites[5].DisplayType = ui.DISPLAY_TYPE_RIGHT
					} else if g.PendIndex == 9 {
						sprites[7].DisplayType = ui.DISPLAY_TYPE_LEFT
						sprites[8].DisplayType = ui.DISPLAY_TYPE_RIGHT
					}

					g.P1Score = g.Game.Pyramid1.Score()
					g.P2Score = g.Game.Pyramid2.Score()
					if g.Game.State == core.IN_PROGRESS {
						if agent := g.Agents[g.Game.Turn%2]; agent != nil {
							g.UIState = WAITING_FOR_OPP_MOVE
							go func() { g.moveChan <- agent.GenerateMove() }()
						}
					} else {
						g.UIState = GAME_OVER
					}
				} else {
					g.DragSprite.ShadowType = 0
					g.AnimationQueue = append(g.AnimationQueue, ui.NewLinearPathAnimator(g.DragSprite, 15,
						ui.Location{X: g.DragSprite.X, Y: g.DragSprite.Y},
						ui.Location{X: DISCARD_X, Y: DISCARD_Y}, ui.EaseOutCubic, nil))
				}
				g.SelectedCard = nil
				g.DragSprite = nil
				g.PendIndex = -1
			}
		}

		if g.Stroke != nil {
			g.Stroke.Update(cx, cy)
			if g.Stroke.Released {
				g.Stroke = nil
			}
		}
	}
}
