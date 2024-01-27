package scene

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/prizelobby/ebitengine-template/res"
	"github.com/prizelobby/ebitengine-template/ui"
	"github.com/tinne26/etxt"
)

const VERSION_STRING = "version 1.0.0"

const CENTER = 640
const TITLE_Y_CENTER = 100
const PLAYING_Y_CENTER = 400
const CREDITS_Y_CENTER = 500

type MenuScene struct {
	BaseScene
	AudioContext *audio.Context
	Sound        []byte
	BgmPlayer    *audio.Player
}

func NewMenuScene(audioContext *audio.Context) *MenuScene {
	b := res.DecodeWavToBytes(audioContext, "dice_03.wav")
	bgm := res.OggToStream(audioContext, "anttisinstrumentals+littleguitar.ogg")
	loop := audio.NewInfiniteLoop(bgm, bgm.Length())
	player, err := audioContext.NewPlayer(loop)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &MenuScene{
		AudioContext: audioContext,
		Sound:        b,
		BgmPlayer:    player,
	}
}

func (m *MenuScene) OnSwitch() {
	m.BgmPlayer.SetPosition(0)
	m.BgmPlayer.Play()
}

func (m *MenuScene) Update() {
	cursorX, cursorY := ui.AdjustedCursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if math.Abs(cursorX-CENTER) < 100 && math.Abs(cursorY-PLAYING_Y_CENTER) < 50 {
			m.SceneManager.SwitchToScene("playing")
			player := m.AudioContext.NewPlayerFromBytes(m.Sound)
			player.Play()
			m.BgmPlayer.Close()
		}

		if math.Abs(cursorX-CENTER) < 100 && math.Abs(cursorY-CREDITS_Y_CENTER) < 50 {
			m.SceneManager.SwitchToScene("credits")
			m.BgmPlayer.Close()
		}
	}
}

func (m *MenuScene) Draw(scaledScreen *ui.ScaledScreen) {
	scaledScreen.DrawTextCenteredAt("EBITENGINE TEMPLATE", 48.0, CENTER, TITLE_Y_CENTER, color.White)
	scaledScreen.DrawTextCenteredAt("Play", 32.0, CENTER, PLAYING_Y_CENTER, color.White)
	scaledScreen.DrawTextCenteredAt("Credits", 32.0, CENTER, CREDITS_Y_CENTER, color.White)

	scaledScreen.DrawTextWithAlign(VERSION_STRING, 16.0, 1280-10, 720-10, color.White, etxt.Bottom, etxt.Right)
}
