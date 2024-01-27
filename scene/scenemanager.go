package scene

import (
	"errors"

	"github.com/prizelobby/ebitengine-template/ui"
)

type Scene interface {
	Update()
	Draw(screen *ui.ScaledScreen)
	OnSwitch()
	SetSceneManager(sm *SceneManager)
}

type SceneManager struct {
	CurrentScene Scene
	SceneDict    map[string]Scene
}

func NewSceneManager() *SceneManager {
	return &SceneManager{
		CurrentScene: &BaseScene{},
		SceneDict:    make(map[string]Scene),
	}
}

func (s *SceneManager) AddScene(name string, scene Scene) {
	scene.SetSceneManager(s)
	s.SceneDict[name] = scene
}

func (s *SceneManager) SwitchToScene(name string) error {
	if nextScene, ok := s.SceneDict[name]; ok {
		s.CurrentScene.OnSwitch()
		s.CurrentScene = nextScene
		return nil
	}
	return errors.New("Scene not found in dict")
}

func (s *SceneManager) Update() {
	s.CurrentScene.Update()
}

func (s *SceneManager) Draw(screen *ui.ScaledScreen) {
	s.CurrentScene.Draw(screen)
}
