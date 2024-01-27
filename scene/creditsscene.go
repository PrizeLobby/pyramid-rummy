package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/prizelobby/ebitengine-template/ui"
)

type CreditsScene struct {
	BaseScene
}

func NewCreditsScene() *CreditsScene {
	return &CreditsScene{}
}

func (c *CreditsScene) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		c.SceneManager.SwitchToScene("menu")
	}
}

func (c *CreditsScene) Draw(screen *ui.ScaledScreen) {
	screen.DrawTextCenteredAt("Credits", 48, 480, 50, color.White)
	screen.DrawTextCenteredAt("click anywhere to return", 16, 480, 450, color.White)
}

func (c *CreditsScene) OnSwitch() {
}
