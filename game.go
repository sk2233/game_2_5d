package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	Game *Game0
)

type Game0 struct {
	Background *Background
}

func NewGame() *Game0 {
	Game = &Game0{Background: NewBackground("resources/bg")}
	return Game
}

func (m *Game0) Update() error {
	Camera.Update()
	m.Background.Update()
	Player.Update()
	RangeManager.Update()
	return nil
}

func (m *Game0) Draw(screen *ebiten.Image) {
	m.Background.Draw(screen)
	Player.Draw(screen)
}

func (m *Game0) Layout(w, h int) (int, int) {
	return w / 4, h / 4
}
