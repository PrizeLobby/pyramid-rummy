package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/prizelobby/pyramid-rummy/core"
	"github.com/prizelobby/pyramid-rummy/res"
	"github.com/prizelobby/pyramid-rummy/scene"
	"github.com/prizelobby/pyramid-rummy/ui"
	einput "github.com/quasilyte/ebitengine-input"
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
	inputSystem  einput.System
}

func (g *EbitenGame) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) && runtime.GOOS != "js" {
		os.Exit(0)
	}

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
	scale := ebiten.Monitor().DeviceScaleFactor()
	canvasWidth := GAME_WIDTH * scale
	canvasHeight := GAME_HEIGHT * scale
	return canvasWidth, canvasHeight
}

func main() {

	args := os.Args[1:]

	if len(args) > 0 {
		if args[0] == "agenttest" {
			iterations := 10
			if len(args) > 1 {
				var err error
				iterations, err = strconv.Atoi(args[1])
				if err != nil {
					fmt.Println("unable to parse iterations, defaulting to 10")
				}
			}
			p1Wins := 0
			p2Wins := 0
			p1TotalScore := 0
			p2TotalScore := 0
			draws := 0
			for i := range iterations {
				game := core.NewGame()
				a1 := core.NewSampleAgent(0)
				a2 := core.NewSampleAgent(1)
				a2.Strategy = 1
				var currentAgent core.GameAgent = a1
				var otherAgent core.GameAgent = a2
				for game.State == core.IN_PROGRESS {
					currentAgent.SetVisibleCard(game.TopDiscard())
					m := currentAgent.GenerateMove()
					if m.EventType == core.DRAW_CARDS {
						c := game.DrawCard()
						currentAgent.RevealCard(c)
						otherAgent.RevealCard(c)
					} else if m.EventType == core.PLAY_CARD {
						t := m.Target
						game.PlayCard(t)
						otherAgent.AcceptMove(game.TopDiscard(), t)
						currentAgent, otherAgent = otherAgent, currentAgent
					}
				}
				fmt.Printf("Game %d\n", i)
				if game.State == core.P1_WIN {
					fmt.Println("p1 win")
					p1Wins += 1
				} else if game.State == core.P2_WIN {
					fmt.Println("p2 win")
					p2Wins += 1
				} else {
					fmt.Println("draw")
					draws += 1
				}
				p1TotalScore += game.Pyramid1.Score()
				p2TotalScore += game.Pyramid2.Score()
				fmt.Printf("Score %d - %d\n", game.Pyramid1.Score(), game.Pyramid2.Score())
			}
			fmt.Printf("Results %d %d %d\n", p1Wins, p2Wins, draws)
			fmt.Printf("Avg scores %.2f %.2f\n", float64(p1TotalScore)/float64(iterations), float64(p2TotalScore)/float64(iterations))
		}
		os.Exit(0)
	}

	audioContext := audio.NewContext(SAMPLE_RATE)
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
	g.inputSystem.Init(einput.SystemConfig{
		DevicesEnabled: einput.AnyDevice,
	})

	sm := scene.NewSceneManager()
	menuScene := scene.NewMenuScene(audioContext)
	creditsScene := scene.NewCreditsScene()

	sm.AddScene("menu", menuScene)
	sm.AddScene("credits", creditsScene)

	g.SceneManager = sm
	sm.SwitchToScene("menu")

	ebiten.SetWindowSize(GAME_WIDTH, GAME_HEIGHT)
	ebiten.SetWindowTitle("Rummy Pyramid")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
