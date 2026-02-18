package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	Player = NewPlayer()
)

const (
	PlayerMove = iota
	PlayerJump
	PlayerFly
)

type Player0 struct {
	Pos     Vec
	Dir     float64
	Anim    *Animator
	Speed   float64
	State   int
	Actions []func()
}

func (p *Player0) Update() {
	p.Anim.Update()
	p.Actions[p.State]()
}

func (p *Player0) Move() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.State = PlayerJump
		p.Speed = 12
		return
	}
	speed := Vec{X: GetAxis(ebiten.KeyD, ebiten.KeyA), Y: GetAxis(ebiten.KeyS, ebiten.KeyW)}.Scale(Speed)
	if speed.X == 0 && speed.Y == 0 {
		return
	}
	if speed.X != 0 {
		p.Dir = Sign(speed.X)
	}
	pos := p.Pos.Add(speed)
	pos.Z = RangeManager.GetZ(pos)
	if math.Abs(pos.Z-p.Pos.Z) <= Speed {
		p.Pos = pos
	} else if pos.Z < p.Pos.Z {
		p.Pos.X, p.Pos.Y = pos.X, pos.Y
		p.Speed = 0
		p.State = PlayerJump
	}
}

func (p *Player0) Jump() {
	p.Speed -= G
	p.Pos.Z += p.Speed
	z := RangeManager.GetZ(p.Pos)
	if p.Pos.Z <= z {
		p.Pos.Z = z
		p.State = PlayerMove
		return
	}
	speed := Vec{X: GetAxis(ebiten.KeyD, ebiten.KeyA), Y: GetAxis(ebiten.KeyS, ebiten.KeyW)}.Scale(Speed)
	if speed.X == 0 && speed.Y == 0 {
		return
	}
	if speed.X != 0 {
		p.Dir = Sign(speed.X)
	}
	pos := p.Pos.Add(speed)
	pos.Z = RangeManager.GetZ(pos)
	if p.Pos.Z >= pos.Z {
		p.Pos.X, p.Pos.Y = pos.X, pos.Y
	}
}

func (p *Player0) Fly() {
	p.Speed -= G
	p.Pos.Z += p.Speed
	z := RangeManager.GetZ(p.Pos)
	if p.Pos.Z <= z {
		p.Pos.Z = z
		p.State = PlayerMove
		return
	}
	speed := Vec{X: 3}
	pos := p.Pos.Add(speed)
	pos.Z = RangeManager.GetZ(pos)
	if p.Pos.Z >= pos.Z {
		p.Pos.X, p.Pos.Y = pos.X, pos.Y
	}
}

func (p *Player0) Draw(screen *ebiten.Image) {
	p.Anim.Draw(screen, ToVec2(p.Pos), p.Dir)
}

func (p *Player0) Hurt() {
	p.Speed = 12
	p.State = PlayerFly
}

func NewPlayer() *Player0 {
	res := &Player0{Anim: NewAnimator(), State: PlayerMove, Dir: 1}
	res.Anim.AddAnim(NewAnim("resources/player").SetPivot(PivotBottom))
	res.Pos = Vec{X: 50, Y: 100}
	res.Pos.Z = RangeManager.GetZ(res.Pos)
	res.Actions = make([]func(), 3)
	res.Actions[PlayerMove] = res.Move
	res.Actions[PlayerJump] = res.Jump
	res.Actions[PlayerFly] = res.Fly
	return res
}
