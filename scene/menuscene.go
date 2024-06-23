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
	"github.com/prizelobby/ebitengine-template/util"
)

const VERSION_STRING = "version 1.0.0"

const CENTER = 640
const TITLE_Y_CENTER = 180
const CHOICE_HEADER_Y = 280
const PLAYING_Y_CENTER = 450
const CREDITS_Y_CENTER = 550

type MenuScene struct {
	BaseScene
	AudioContext *audio.Context
	Sound        []byte
	BgmPlayer    *audio.Player

	P1Choice int
	P2Choice int
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
		P2Choice:     1,
	}
}

func (m *MenuScene) OnSwitch() {
	m.BgmPlayer.SetPosition(0)
	m.BgmPlayer.Play()
}

func (m *MenuScene) Update() {
	cx, cy := ui.AdjustedCursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if math.Abs(cx-CENTER) < 100 && math.Abs(cy-PLAYING_Y_CENTER) < 50 {
			m.SceneManager.SwitchToScene("game")
			player := m.AudioContext.NewPlayerFromBytes(m.Sound)
			player.Play()
			m.BgmPlayer.Close()
		}
		if util.XYinRect(cx, cy, CENTER-100-48, CHOICE_HEADER_Y+40-20, 48*2, 20*2) {
			m.P1Choice = 0
		} else if util.XYinRect(cx, cy, CENTER-100-48, CHOICE_HEADER_Y+80-20, 48*2, 20*2) {
			m.P1Choice = 1
		} else if util.XYinRect(cx, cy, CENTER+100-48, CHOICE_HEADER_Y+40-20, 48*2, 20*2) {
			m.P2Choice = 0
		} else if util.XYinRect(cx, cy, CENTER+100-48, CHOICE_HEADER_Y+80-20, 48*2, 20*2) {
			m.P2Choice = 1
		}

		/*
			if math.Abs(cursorX-CENTER) < 100 && math.Abs(cursorY-CREDITS_Y_CENTER) < 50 {
				m.SceneManager.SwitchToScene("credits")
				m.BgmPlayer.Close()
			}*/
	}
}

func (m *MenuScene) Draw(screen *ui.ScaledScreen) {
	screen.Screen.Fill(color.RGBA{0x5b, 0xa1, 0x2d, 0xff})

	screen.DrawTextCenteredAt("Rummy Pyramid", 56.0, CENTER, TITLE_Y_CENTER, color.White)
	screen.DrawTextCenteredAt("Play", 48.0, CENTER, PLAYING_Y_CENTER, color.White)
	screen.DrawTextCenteredAt("Rules", 48.0, CENTER, CREDITS_Y_CENTER, color.White)

	screen.DrawTextCenteredAt("Player 1", 32.0, CENTER-100, CHOICE_HEADER_Y, color.White)
	screen.DrawTextCenteredAt("Human", 24.0, CENTER-100, CHOICE_HEADER_Y+40, color.White)
	screen.DrawTextCenteredAt("Computer", 24.0, CENTER-100, CHOICE_HEADER_Y+80, color.White)
	screen.DrawTextCenteredAt("Player 2", 32.0, CENTER+100, CHOICE_HEADER_Y, color.White)
	screen.DrawTextCenteredAt("Human", 24.0, CENTER+100, CHOICE_HEADER_Y+40, color.White)
	screen.DrawTextCenteredAt("Computer", 24.0, CENTER+100, CHOICE_HEADER_Y+80, color.White)

	if m.P1Choice == 0 {
		screen.DrawCircle(CENTER-100-48, CHOICE_HEADER_Y+40, 4, color.White)
		screen.DrawCircle(CENTER-100+48, CHOICE_HEADER_Y+40, 4, color.White)
	} else {
		screen.DrawCircle(CENTER-100-60, CHOICE_HEADER_Y+80, 4, color.White)
		screen.DrawCircle(CENTER-100+60, CHOICE_HEADER_Y+80, 4, color.White)
	}

	if m.P2Choice == 0 {
		screen.DrawCircle(CENTER+100-48, CHOICE_HEADER_Y+40, 4, color.White)
		screen.DrawCircle(CENTER+100+48, CHOICE_HEADER_Y+40, 4, color.White)
	} else {
		screen.DrawCircle(CENTER+100-60, CHOICE_HEADER_Y+80, 4, color.White)
		screen.DrawCircle(CENTER+100+60, CHOICE_HEADER_Y+80, 4, color.White)
	}

	//scaledScreen.DrawTextCenteredAt("Credits", 32.0, CENTER, CREDITS_Y_CENTER, color.White)

	//scaledScreen.DrawTextWithAlign(VERSION_STRING, 16.0, 1280-10, 720-10, color.White, etxt.Bottom, etxt.Right)
}
