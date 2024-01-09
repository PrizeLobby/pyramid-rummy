package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/prizelobby/ebitengine-template/core"
	"github.com/prizelobby/ebitengine-template/res"
	"github.com/prizelobby/ebitengine-template/scene"
	"github.com/prizelobby/ebitengine-template/ui"
	"github.com/tinne26/etxt"
)

const GAME_WIDTH = 1280
const GAME_HEIGHT = 720

const SAMPLE_RATE = 48000

type GameState int

const (
	MENU GameState = iota
	PLAYING
	CREDITS
)

type EbitenGame struct {
	ScaledScreen *ui.ScaledScreen
	gameState    GameState
	SceneManager *scene.SceneManager
}

func (g *EbitenGame) Update() error {
	g.SceneManager.Update()
	return nil
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
	g.ScaledScreen.SetTarget(screen)

	//msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	//g.ScaledScreen.DebugPrint(msg)

	g.SceneManager.Draw(g.ScaledScreen)
}

func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	panic("use Ebitengine >=v2.5.0")
}

func (g *EbitenGame) LayoutF(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	scale := ebiten.DeviceScaleFactor()
	canvasWidth := GAME_WIDTH * scale
	canvasHeight := GAME_HEIGHT * scale
	return canvasWidth, canvasHeight
}

func main() {
	audioContext := audio.NewContext(SAMPLE_RATE)
	game := core.NewGame()

	// create a new text renderer and configure it
	txtRenderer := etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	txtRenderer.SetCacheHandler(glyphsCache.NewHandler())
	txtRenderer.SetFont(res.GetFont("Roboto-Medium"))
	txtRenderer.SetAlign(etxt.YCenter, etxt.XCenter)
	txtRenderer.SetSizePx(64)

	scaledScreen := ui.NewScaledScreen(txtRenderer)

	g := &EbitenGame{
		ScaledScreen: scaledScreen,
		gameState:    MENU,
	}
	sm := scene.NewSceneManager()
	menuScene := scene.NewMenuScene(audioContext)
	creditsScene := scene.NewCreditsScene()
	gameScene := scene.NewGameScene(game)
	sm.AddScene("menu", menuScene)
	sm.AddScene("credits", creditsScene)
	sm.AddScene("playing", gameScene)
	g.SceneManager = sm
	sm.SwitchToScene("menu")

	ebiten.SetWindowSize(GAME_WIDTH, GAME_HEIGHT)
	ebiten.SetWindowTitle("EBITENGINE TEMPLATE")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
