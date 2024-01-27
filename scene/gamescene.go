package scene

import (
	"image/color"

	"github.com/prizelobby/ebitengine-template/core"
	"github.com/prizelobby/ebitengine-template/ui"
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
	Game *core.Game
}

func NewGameScene(game *core.Game) *GameScene {
	return &GameScene{
		Game: game,
	}
}

func (g *GameScene) Draw(screen *ui.ScaledScreen) {
	screen.DrawText("Game Scene", 16, 0, 0, color.White)
}

func (g *GameScene) Update() {

}

func (g *GameScene) OnSwitch() {
}

func (g *GameScene) SetSceneManager(sm *SceneManager) {
	g.SceneManager = sm
}
