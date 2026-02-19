package main

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type Anim struct {
	Name      string
	Pivot     Vec
	Images    []*ebiten.Image
	Option    *ebiten.DrawImageOptions
	StopFrame int
	Speed     float64
	Current   float64
}

func (a *Anim) Update() bool {
	if int(a.Current) == a.StopFrame {
		return false
	}
	a.Current += a.Speed
	if a.Current >= float64(len(a.Images)) {
		a.Current = 0
		return true
	}
	return false
}

func NewAnim(path string) *Anim {
	idx := strings.LastIndex(path, "/")
	return &Anim{Name: strings.TrimSpace(path[idx+1:]), Images: LoadImages(path),
		Option: &ebiten.DrawImageOptions{}, StopFrame: -1}
}

func (a *Anim) SetFps(fps float64) *Anim {
	a.Speed = fps / Fps
	return a
}

var (
	PivotCenter = Vec{X: 0.5, Y: 0.5}
	PivotBottom = Vec{X: 0.5, Y: 1}
)

func (a *Anim) SetPivot(vec Vec) *Anim {
	bound := a.Images[0].Bounds()
	a.Pivot = Vec{X: vec.X * float64(bound.Dx()), Y: vec.Y * float64(bound.Dy())}
	return a
}

func (a *Anim) Draw(screen *ebiten.Image, vec Vec, dir float64) { // 只使用前 2 维
	vec.X -= Camera.X
	a.Option.GeoM.Reset()
	a.Option.GeoM.Translate(-a.Pivot.X, -a.Pivot.Y)
	a.Option.GeoM.Scale(dir, 1)
	a.Option.GeoM.Translate(vec.X, vec.Y)
	screen.DrawImage(a.Images[int(a.Current)], a.Option)
}

func (a *Anim) Reset() {
	a.Current = 0
}

func (a *Anim) FreezeFrame(stopFrame int) *Anim {
	a.StopFrame = stopFrame
	return a
}

func (a *Anim) SetFrame(frame int) {
	a.Current = float64(frame)
}

type Animator struct {
	Anims   map[string]*Anim
	Anim    *Anim
	AnimEnd func(*Anim)
}

func (a *Animator) AddAnim(anim *Anim) {
	a.Anims[anim.Name] = anim
	a.Anim = anim
}

func (a *Animator) Update() {
	if a.Anim.Update() && a.AnimEnd != nil {
		a.AnimEnd(a.Anim)
	}
}

func (a *Animator) Draw(screen *ebiten.Image, pos Vec, dir float64) {
	a.Anim.Draw(screen, pos, dir)
}

func (a *Animator) Play(name string) {
	a.Anim = a.Anims[name]
}

func (a *Animator) RePlay(name string) {
	a.Anim = a.Anims[name]
	a.Anim.Reset()
}

func (a *Animator) SetFrame(frame int) {
	a.Anim.SetFrame(frame)
}

func (a *Animator) SetAnimEnd(animEnd func(*Anim)) {
	a.AnimEnd = animEnd
}

func NewAnimator() *Animator {
	return &Animator{Anims: make(map[string]*Anim)}
}
