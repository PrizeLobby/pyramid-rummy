package scene

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/prizelobby/pyramid-rummy/res"
	"github.com/prizelobby/pyramid-rummy/ui"
	"github.com/prizelobby/pyramid-rummy/util"
)

const VERSION_STRING = "version 1.0.0"

const CENTER = 640
const TITLE_Y_CENTER = 180
const CHOICE_HEADER_Y = 280
const PLAYING_Y_CENTER = 450
const RULES_Y_CENTER = 550

type MenuScene struct {
	BaseScene
	AudioContext *audio.Context
	Sound        []byte

	P0Choice int
	P1Choice int

	Rules        *ui.RulesComponent
	ShowingRules bool
}

func NewMenuScene(audioContext *audio.Context) *MenuScene {
	b := res.DecodeWavToBytes(audioContext, "dice_03.wav")

	return &MenuScene{
		AudioContext: audioContext,
		Sound:        b,
		P1Choice:     1,

		Rules: ui.NewRulesComponent(),
	}
}

func (m *MenuScene) OnSwitch() {
}

func (m *MenuScene) Update() {
	if m.ShowingRules {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			m.ShowingRules = false
		}
		return
	}

	cx, cy := ui.AdjustedCursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		if math.Abs(cx-CENTER) < 100 && math.Abs(cy-PLAYING_Y_CENTER) < 50 {
			gs := NewGameScene(m.P0Choice, m.P1Choice, m.AudioContext)
			m.SceneManager.AddScene("game", gs)
			m.SceneManager.SwitchToScene("game")
			player := m.AudioContext.NewPlayerFromBytes(m.Sound)
			player.Play()
		}
		if util.XYinRect(cx, cy, CENTER-100-48, CHOICE_HEADER_Y+40-20, 48*2, 20*2) {
			m.P0Choice = 0
		} else if util.XYinRect(cx, cy, CENTER-100-48, CHOICE_HEADER_Y+80-20, 48*2, 20*2) {
			m.P0Choice = 1
		} else if util.XYinRect(cx, cy, CENTER+100-48, CHOICE_HEADER_Y+40-20, 48*2, 20*2) {
			m.P1Choice = 0
		} else if util.XYinRect(cx, cy, CENTER+100-48, CHOICE_HEADER_Y+80-20, 48*2, 20*2) {
			m.P1Choice = 1
		} else if util.XYinRect(cx, cy, CENTER-48, RULES_Y_CENTER-20, 48*2, 20*2) {
			m.ShowingRules = true
		}

		/*
			if math.Abs(cursorX-CENTER) < 100 && math.Abs(cursorY-CREDITS_Y_CENTER) < 50 {
				m.SceneManager.SwitchToScene("credits")
				m.BgmPlayer.Close()
			}*/
	}
}

func (m *MenuScene) Draw(screen *ui.ScaledScreen) {
	screen.Screen.Fill(color.RGBA{0x44, 0x5c, 0x47, 0xff})

	if m.ShowingRules {
		m.Rules.Draw(screen)
		return
	}

	screen.DrawTextCenteredAt("Rummy Pyramid", 56.0, CENTER, TITLE_Y_CENTER, color.White)
	screen.DrawTextCenteredAt("Play", 48.0, CENTER, PLAYING_Y_CENTER, color.White)
	screen.DrawTextCenteredAt("Rules", 48.0, CENTER, RULES_Y_CENTER, color.White)

	screen.DrawTextCenteredAt("Player 1", 32.0, CENTER-100, CHOICE_HEADER_Y, color.White)
	screen.DrawTextCenteredAt("Human", 24.0, CENTER-100, CHOICE_HEADER_Y+40, color.White)
	screen.DrawTextCenteredAt("Computer", 24.0, CENTER-100, CHOICE_HEADER_Y+80, color.White)
	screen.DrawTextCenteredAt("Player 2", 32.0, CENTER+100, CHOICE_HEADER_Y, color.White)
	screen.DrawTextCenteredAt("Human", 24.0, CENTER+100, CHOICE_HEADER_Y+40, color.White)
	screen.DrawTextCenteredAt("Computer", 24.0, CENTER+100, CHOICE_HEADER_Y+80, color.White)

	if m.P0Choice == 0 {
		screen.DrawCircle(CENTER-100-48, CHOICE_HEADER_Y+40, 4, color.White)
		screen.DrawCircle(CENTER-100+48, CHOICE_HEADER_Y+40, 4, color.White)
	} else {
		screen.DrawCircle(CENTER-100-60, CHOICE_HEADER_Y+80, 4, color.White)
		screen.DrawCircle(CENTER-100+60, CHOICE_HEADER_Y+80, 4, color.White)
	}

	if m.P1Choice == 0 {
		screen.DrawCircle(CENTER+100-48, CHOICE_HEADER_Y+40, 4, color.White)
		screen.DrawCircle(CENTER+100+48, CHOICE_HEADER_Y+40, 4, color.White)
	} else {
		screen.DrawCircle(CENTER+100-60, CHOICE_HEADER_Y+80, 4, color.White)
		screen.DrawCircle(CENTER+100+60, CHOICE_HEADER_Y+80, 4, color.White)
	}

	//scaledScreen.DrawTextCenteredAt("Credits", 32.0, CENTER, CREDITS_Y_CENTER, color.White)
}
