package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Background struct {
	Anim *Anim
}

func (b *Background) Update() {
	b.Anim.Update()
}

func (b *Background) Draw(screen *ebiten.Image) {
	b.Anim.Draw(screen, Vec{}, 1)
}

func NewBackground(path string) *Background {
	return &Background{Anim: NewAnim(path).SetFps(9)}
}
