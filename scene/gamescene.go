package scene

import (
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/prizelobby/pyramid-rummy/core"
	"github.com/prizelobby/pyramid-rummy/res"
	"github.com/prizelobby/pyramid-rummy/ui"
	"github.com/prizelobby/pyramid-rummy/util"
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

	Input       *einput.Handler
	PendIndex   int
	PrevPend    int
	P0Score     int
	P1Score     int
	CurrentTurn int // Used to delay some ui updates because of animations

	DragSprite     *ui.CardSprite
	DiscardSprite  *ui.CardSprite
	SecondSprite   *ui.CardSprite
	P0Spheres      [10]*ui.CardSprite
	P1Spheres      [10]*ui.CardSprite
	MapSmall       *ebiten.Image
	HexMap         *ebiten.Image
	HexMapInactive *ebiten.Image
	BaseTile       *ebiten.Image
	HoverTile      *ebiten.Image
	OutlineTile    *ebiten.Image
	Shadow         *ebiten.Image

	Stroke *ui.Stroke

	ActiveAnimation ui.Anim
	AnimationQueue  []ui.Anim

	HelpText string
}

func NewGameScene(p0, p1 int) *GameScene {
	game := core.NewGame()

	agents := [2]core.GameAgent{nil, nil}
	if p0 == 1 {
		agents[0] = core.NewAgent(0)
	}
	if p1 == 1 {
		agents[1] = core.NewAgent(1)
	}

	return &GameScene{
		Game:           game,
		HexMap:         res.GetImage("hexmap"),
		HexMapInactive: res.GetImage("hexmapdeselected"),
		BaseTile:       res.GetImage("basetile"),
		HoverTile:      res.GetImage("hexoutlinegreen"),
		OutlineTile:    res.GetImage("hexoutlinebroken"),
		MapSmall:       res.GetImage("circlemapsmall"),
		Shadow:         res.GetImage("shadow"),
		moveChan:       make(chan core.AgentEvent, 1),
		PendIndex:      -1,
		HelpText:       "Click the deck to reveal a card.",
		Agents:         agents,
	}
}

const HELPTEXT_Y = 180
const TURN_TEXT_Y = 90

const DECK_BUTTON_X = 640 - 60 - ui.TILE_SIZE_X
const DECK_BUTTON_Y = 320
const DECK_BUTTON_W = 120
const DECK_BUTTON_H = 146

const DISCARD_X = 640 + 60
const DISCARD_Y = 320

const P0StartX float64 = 80
const P0StartY float64 = 290
const P1StartX float64 = 1280 - 356 - P0StartX
const P1StartY float64 = 290

var P0XLocs [10]float64 = [10]float64{
	P0StartX + ui.TILE_X_OFFSET, P0StartX + ui.TILE_X_OFFSET/2, P0StartX + ui.TILE_X_OFFSET*1.5, P0StartX, P0StartX + ui.TILE_X_OFFSET, P0StartX + 2*ui.TILE_X_OFFSET,
	P0StartX + ui.TILE_X_OFFSET, P0StartX + ui.TILE_X_OFFSET/2, P0StartX + ui.TILE_X_OFFSET*1.5, P0StartX + ui.TILE_X_OFFSET,
}
var P0YLocs [10]float64 = [10]float64{
	P0StartY, P0StartY + ui.TILE_Y_OFFSET, P0StartY + ui.TILE_Y_OFFSET, P0StartY + 2*ui.TILE_Y_OFFSET, P0StartY + 2*ui.TILE_Y_OFFSET, P0StartY + 2*ui.TILE_Y_OFFSET,
	P0StartY + math.Floor(ui.TILE_Y_OFFSET/2), P0StartY + math.Floor(ui.TILE_Y_OFFSET*1.5), P0StartY + math.Floor(ui.TILE_Y_OFFSET*1.5),
	P0StartY + ui.TILE_Y_OFFSET - 1,
}

var P1XLocs [10]float64 = [10]float64{
	P1StartX + ui.TILE_X_OFFSET, P1StartX + ui.TILE_X_OFFSET/2, P1StartX + ui.TILE_X_OFFSET*1.5, P1StartX, P1StartX + ui.TILE_X_OFFSET, P1StartX + 2*ui.TILE_X_OFFSET,
	P1StartX + ui.TILE_X_OFFSET, P1StartX + ui.TILE_X_OFFSET/2, P1StartX + ui.TILE_X_OFFSET*1.5, P1StartX + ui.TILE_X_OFFSET,
}
var P1YLocs [10]float64 = [10]float64{
	P0StartY, P0StartY + ui.TILE_Y_OFFSET, P0StartY + ui.TILE_Y_OFFSET, P0StartY + 2*ui.TILE_Y_OFFSET, P0StartY + 2*ui.TILE_Y_OFFSET, P0StartY + 2*ui.TILE_Y_OFFSET,
	P0StartY + math.Floor(ui.TILE_Y_OFFSET/2), P0StartY + math.Floor(ui.TILE_Y_OFFSET*1.5), P0StartY + math.Floor(ui.TILE_Y_OFFSET*1.5),
	P0StartY + ui.TILE_Y_OFFSET - 1,
}

func (g *GameScene) Draw(screen *ui.ScaledScreen) {
	screen.Screen.Fill(color.RGBA{0x44, 0x5c, 0x47, 0xff})

	screen.DrawTextCenteredAt(g.HelpText, 32, 640, HELPTEXT_Y, color.White)
	if g.UIState != GAME_OVER {
		screen.DrawTextCenteredAt("Player "+strconv.Itoa(g.CurrentTurn+1)+"'s turn", 48, 640, TURN_TEXT_Y, color.White)
	} else {
		if g.Game.State == core.P1_WIN {
			screen.DrawTextCenteredAt("P1 Wins", 48, 640, TURN_TEXT_Y, color.White)
		} else if g.Game.State == core.P2_WIN {
			screen.DrawTextCenteredAt("P2 Wins", 48, 640, TURN_TEXT_Y, color.White)
		} else if g.Game.State == core.DRAW {
			screen.DrawTextCenteredAt("Draw", 48, 640, TURN_TEXT_Y, color.White)
		}
	}

	deckOpts := &ebiten.DrawImageOptions{}
	deckOpts.GeoM.Translate(DECK_BUTTON_X, DECK_BUTTON_Y)
	screen.DrawImage(g.BaseTile, deckOpts)
	deckShadowOpts := &ebiten.DrawImageOptions{}
	deckShadowOpts.GeoM.Translate(DECK_BUTTON_X, DECK_BUTTON_Y)
	screen.DrawImage(g.Shadow, deckShadowOpts)
	//screen.DrawTextCenteredAt(strconv.Itoa(len(g.Game.Deck))+"\nCards Left", 20, DECK_BUTTON_X+ui.TILE_X_OFFSET/2, DECK_BUTTON_Y+60, color.Black)

	screen.DrawTextCenteredAt("Score: "+strconv.Itoa(g.P0Score), 36, P0StartX+ui.TILE_X_OFFSET*1.5, P0StartY-40, color.White)
	screen.DrawTextCenteredAt("Score: "+strconv.Itoa(g.P1Score), 36, P1StartX+ui.TILE_X_OFFSET*1.5, P1StartY-40, color.White)

	pOpts := &ebiten.DrawImageOptions{}
	pOpts.GeoM.Translate(P0StartX, P0StartY)
	pOpts2 := &ebiten.DrawImageOptions{}
	pOpts2.GeoM.Translate(P1StartX, P0StartY)
	if g.CurrentTurn == 0 {
		screen.DrawImage(g.HexMap, pOpts)
		screen.DrawImage(g.HexMapInactive, pOpts2)
	} else {
		screen.DrawImage(g.HexMapInactive, pOpts)
		screen.DrawImage(g.HexMap, pOpts2)
	}

	screen.DrawTextCenteredAt("Revealed:", 30, DISCARD_X+ui.TILE_X_OFFSET/2, DISCARD_Y-50, color.White)
	screen.DrawTextCenteredAt("Deck:", 30, DECK_BUTTON_X+ui.TILE_X_OFFSET/2, DECK_BUTTON_Y-50, color.White)
	if g.Agents[g.Game.Turn%2] == nil && g.CurrentTurn == g.Game.Turn%2 {
		plural := "s"
		if g.Game.DrawsLeft == 1 {
			plural = ""
		}
		screen.DrawTextCenteredAt(strconv.Itoa(g.Game.DrawsLeft)+" draw"+plural+" left", 24, DECK_BUTTON_X+ui.TILE_X_OFFSET/2, DECK_BUTTON_Y-20, color.White)
	}

	dOpts := &ebiten.DrawImageOptions{}
	dOpts.GeoM.Translate(DISCARD_X, DISCARD_Y+ui.TILE_HEIGHT)
	screen.DrawImage(g.OutlineTile, dOpts)
	screen.DrawTextCenteredAt("No cards\nin stack", 18, DISCARD_X+ui.TILE_X_OFFSET/2, DISCARD_Y+70, color.White)

	if g.SecondSprite != nil {
		g.SecondSprite.Draw(screen)
	}
	if g.DiscardSprite != nil {
		g.DiscardSprite.Draw(screen)
	}

	if g.PendIndex != -1 && g.PendIndex <= 5 {
		opt := &ebiten.DrawImageOptions{}
		_, x, y := g.PyramidXYForTurn(g.PendIndex)
		opt.GeoM.Translate(x, y)
		screen.DrawImage(g.HoverTile, opt)
	}

	for i, s := range g.P0Spheres {
		if s != nil {
			s.Draw(screen)
		}
		if i == 5 {
			if g.CurrentTurn == 0 && g.Agents[0] == nil {
				for j := 6; j < 9; j++ {
					if g.Game.Pyramid1.CanPlace(j) {
						outlineOpt := &ebiten.DrawImageOptions{}
						outlineOpt.GeoM.Translate(P0XLocs[j], P0YLocs[j])
						screen.DrawImage(g.OutlineTile, outlineOpt)
					}
				}
				if g.PendIndex > 5 {
					opt := &ebiten.DrawImageOptions{}
					opt.GeoM.Translate(P0XLocs[g.PendIndex], P0YLocs[g.PendIndex])
					screen.DrawImage(g.HoverTile, opt)
				}
			}
		}
	}
	for i, s := range g.P1Spheres {
		if s != nil {
			s.Draw(screen)
		}
		if i == 5 {
			if g.CurrentTurn == 1 && g.Agents[1] == nil {
				for j := 6; j < 9; j++ {
					if g.Game.Pyramid2.CanPlace(j) {
						outlineOpt := &ebiten.DrawImageOptions{}
						outlineOpt.GeoM.Translate(P1XLocs[j], P1YLocs[j])
						screen.DrawImage(g.OutlineTile, outlineOpt)
					}
				}
				if g.PendIndex > 5 {
					opt := &ebiten.DrawImageOptions{}
					opt.GeoM.Translate(P1XLocs[g.PendIndex], P1YLocs[g.PendIndex])
					screen.DrawImage(g.HoverTile, opt)
				}
			}
		}
	}

	if g.Game.Pyramid1.CanPlace(9) && g.CurrentTurn == 0 && g.Agents[0] == nil {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(P0XLocs[9], P0YLocs[9])
		screen.DrawImage(g.OutlineTile, opt)
	}
	if g.Game.Pyramid2.CanPlace(9) && g.CurrentTurn == 1 && g.Agents[1] == nil {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(P1XLocs[9], P1YLocs[9])
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

	if g.UIState == GAME_OVER {
		screen.DrawUnfilledRect(640-120, 550-20, 240, 40, 2, color.White)
		screen.DrawTextCenteredAt("Return to menu", 32, 640, 550, color.White)
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

	x, y := P0XLocs[i], P0YLocs[i]
	pyramid := g.Game.Pyramid1
	if g.Game.Turn%2 == 1 {
		x, y = P1XLocs[i], P1YLocs[i]
		pyramid = g.Game.Pyramid2
	}
	return pyramid, x, y
}

func (g *GameScene) Update() {
	select {
	case m := <-g.moveChan:
		if m.EventType == core.DRAW_CARDS {
			c := g.Game.DrawCard()
			g.Agents[g.Game.Turn%2].RevealCard(c)
			g.SecondSprite = g.DiscardSprite
			g.DiscardSprite = ui.NewCardSprite(c, DECK_BUTTON_X, DECK_BUTTON_Y)
			g.AnimationQueue = append(g.AnimationQueue, ui.NewBlockingAnim(30), ui.NewLinearPathAnimator(g.DiscardSprite, 35,
				ui.Location{X: DECK_BUTTON_X, Y: DECK_BUTTON_Y},
				ui.Location{X: DISCARD_X, Y: DISCARD_Y}, ui.EaseOutCubic, func() {
					go func() {
						g.moveChan <- g.Agents[g.Game.Turn%2].GenerateMove()
					}()
				}))

		} else if m.EventType == core.PLAY_CARD {
			g.Game.PlayCard(m.Target)
			complete := func() {
				// this code is almost repeated, but its fine for now
				var sprites [10]*ui.CardSprite
				if g.Game.Turn%2 == 0 {
					sprites = g.P1Spheres
				} else {
					sprites = g.P0Spheres
				}
				if m.Target == 6 {
					sprites[1].DisplayType = ui.DISPLAY_TYPE_LEFT
					sprites[2].DisplayType = ui.DISPLAY_TYPE_RIGHT
				} else if m.Target == 7 {
					if sprites[6] != nil {
						sprites[1].DisplayType = ui.DISPLAY_TYPE_LEFT
					}
					sprites[3].DisplayType = ui.DISPLAY_TYPE_LEFT
					if sprites[8] != nil {
						sprites[4].DisplayType = ui.DISPLAY_TYPE_BOTTOM
					} else {
						sprites[4].DisplayType = ui.DISPLAY_TYPE_RIGHT
					}
				} else if m.Target == 8 {
					if sprites[7] != nil {
						sprites[2].DisplayType = ui.DISPLAY_TYPE_RIGHT
					}
					if sprites[7] != nil {
						sprites[4].DisplayType = ui.DISPLAY_TYPE_BOTTOM
					} else {
						sprites[4].DisplayType = ui.DISPLAY_TYPE_LEFT
					}
					sprites[5].DisplayType = ui.DISPLAY_TYPE_RIGHT
				} else if m.Target == 9 {
					sprites[7].DisplayType = ui.DISPLAY_TYPE_LEFT
					sprites[8].DisplayType = ui.DISPLAY_TYPE_RIGHT
				}

				nextCard := g.Game.TopDiscard()
				if nextCard != nil {
					g.DiscardSprite = ui.NewCardSprite(nextCard, DISCARD_X, DISCARD_Y)
				}
				g.P0Score = g.Game.Pyramid1.Score()
				g.P1Score = g.Game.Pyramid2.Score()
				g.CurrentTurn = 1 - g.CurrentTurn
				if g.Game.State == core.IN_PROGRESS {
					if g.Agents[g.Game.Turn%2] == nil {
						g.UIState = WAITING_FOR_PLAYER_MOVE
						if len(g.Game.Discards) > 0 {
							g.HelpText = "Drag the open card to your pyramid or click the deck to reveal a new card."
						} else {
							g.HelpText = "Click the deck to reveal a card."
						}
					} else {
						go func() { g.moveChan <- g.Agents[g.Game.Turn%2].GenerateMove() }()
					}
				} else {
					g.UIState = GAME_OVER
					g.HelpText = "Game Over."
				}
			}
			if g.Game.Turn%2 == 0 {
				g.P1Spheres[m.Target] = g.DiscardSprite // ui.NewCardSprite(card, P2XLocs[m.Target], P2YLocs[m.Target])
				g.DiscardSprite = nil
				g.AnimationQueue = append(g.AnimationQueue, ui.NewBlockingAnim(30), ui.NewLinearPathAnimator(g.P1Spheres[m.Target], 50,
					ui.Location{X: DISCARD_X, Y: DISCARD_Y},
					ui.Location{X: P1XLocs[m.Target], Y: P1YLocs[m.Target] - ui.TILE_HEIGHT}, ui.EaseOutCubic, complete))
			} else {
				g.P0Spheres[m.Target] = g.DiscardSprite // ui.NewCardSprite(card, P2XLocs[m.Target], P2YLocs[m.Target])
				g.DiscardSprite = nil
				g.AnimationQueue = append(g.AnimationQueue, ui.NewBlockingAnim(30), ui.NewLinearPathAnimator(g.P0Spheres[m.Target], 50,
					ui.Location{X: DISCARD_X, Y: DISCARD_Y},
					ui.Location{X: P0XLocs[m.Target], Y: P0YLocs[m.Target] - ui.TILE_HEIGHT}, ui.EaseOutCubic, complete))
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
		}
	}

	if g.UIState == GAME_OVER {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			cx, cy := ui.AdjustedCursorPosition()
			if util.XYinRect(cx, cy, 640-120, 550-20, 240, 40) {
				g.SceneManager.SwitchToScene("menu")
			}
		}
	} else if g.UIState == WAITING_FOR_PLAYER_MOVE {
		g.PrevPend = g.PendIndex

		cx, cy := ui.AdjustedCursorPosition()
		InHex := false
		if g.DragSprite != nil {
			for i := range 10 {
				pyramid, x, y := g.PyramidXYForTurn(i)
				if XYinHexCell(cx, cy, x, y, ui.TILE_SIZE_X, ui.TILE_SIZE_Y-ui.TILE_HEIGHT, ui.TILE_TIP_HEIGHT) && pyramid.CanPlace(i) {
					g.PendIndex = i
					if g.PendIndex != g.PrevPend {
						if g.Game.Turn%2 == 0 {
							g.P0Score = pyramid.TentativeScoreWithCard(g.DragSprite.Card, g.PendIndex)
						} else {
							g.P1Score = pyramid.TentativeScoreWithCard(g.DragSprite.Card, g.PendIndex)
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
				g.P0Score = g.Game.Pyramid1.Score()
				g.P1Score = g.Game.Pyramid2.Score()
			}
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if g.DiscardSprite != nil && g.DiscardSprite.In(cx, cy) {
				g.Stroke = ui.NewStroke(cx, cy, g.DiscardSprite, 0)
				g.DragSprite = g.DiscardSprite
				g.DragSprite.ShadowType = 1
				g.DragSprite.X = cx - ui.TILE_SIZE_X/2
				g.DragSprite.Y = cy - (ui.TILE_SIZE_Y-ui.TILE_HEIGHT-6)/2
			} else if util.XYinRect(cx, cy, DECK_BUTTON_X, DECK_BUTTON_Y, DECK_BUTTON_W, DECK_BUTTON_H) {
				if g.Game.DrawsLeft == 0 {
					g.HelpText = "You have 0 draws remaining this turn. Drag the open card to your pyramid."
				} else {
					c := g.Game.DrawCard()
					g.SecondSprite = g.DiscardSprite
					g.DiscardSprite = ui.NewCardSprite(c, DECK_BUTTON_X, DECK_BUTTON_Y)
					g.UIState = WAITING_FOR_PLAYER_ANIMIMATION
					g.AnimationQueue = append(g.AnimationQueue, ui.NewLinearPathAnimator(g.DiscardSprite, 25,
						ui.Location{X: DECK_BUTTON_X, Y: DECK_BUTTON_Y},
						ui.Location{X: DISCARD_X, Y: DISCARD_Y}, ui.EaseOutCubic, func() { g.UIState = WAITING_FOR_PLAYER_MOVE }))
					if g.Game.DrawsLeft == 0 {
						g.HelpText = "Drag the open card to your pyramid."
					} else {
						g.HelpText = "Drag the open card to your pyramid or click the deck to reveal a new card."
					}
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
						if len(g.Game.Discards) > 1 {
							g.SecondSprite = ui.NewCardSprite(g.Game.Discards[len(g.Game.Discards)-2], DISCARD_X, DISCARD_Y)
							g.SecondSprite.X = DISCARD_X
							g.SecondSprite.Y = DISCARD_Y
						} else {
							g.SecondSprite = nil
						}
						g.HelpText = "Drag the open card to your pyramid or click the deck to reveal a new card."
					} else {
						g.DiscardSprite = nil
						g.HelpText = "Click the deck to reveal a card."
					}
					g.DragSprite.X = x
					g.DragSprite.Y = y - ui.TILE_HEIGHT
					if g.PendIndex < 6 {
						g.DragSprite.ShadowType = 0
					} else {
						g.DragSprite.ShadowType = 2
					}
					// this conditional is flipped because we increment the turn counter during play card
					var sprites [10]*ui.CardSprite
					if g.Game.Turn%2 == 0 {
						g.P1Spheres[g.PendIndex] = g.DragSprite
						sprites = g.P1Spheres
					} else {
						g.P0Spheres[g.PendIndex] = g.DragSprite
						sprites = g.P0Spheres
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

					g.P0Score = g.Game.Pyramid1.Score()
					g.P1Score = g.Game.Pyramid2.Score()
					g.CurrentTurn = 1 - g.CurrentTurn
					if g.Game.State == core.IN_PROGRESS {
						if agent := g.Agents[g.Game.Turn%2]; agent != nil {
							g.UIState = WAITING_FOR_OPP_MOVE
							g.HelpText = "The computer is thinking..."
							go func() { g.moveChan <- agent.GenerateMove() }()
						}
					} else {
						g.UIState = GAME_OVER
						g.HelpText = "Game Over. Click anywhere to return to main menu."
					}
				} else {
					g.DragSprite.ShadowType = 0
					g.UIState = WAITING_FOR_PLAYER_ANIMIMATION
					g.AnimationQueue = append(g.AnimationQueue, ui.NewLinearPathAnimator(g.DragSprite, 15,
						ui.Location{X: g.DragSprite.X, Y: g.DragSprite.Y},
						ui.Location{X: DISCARD_X, Y: DISCARD_Y}, ui.EaseOutCubic, func() {
							g.UIState = WAITING_FOR_PLAYER_MOVE
						}))
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
