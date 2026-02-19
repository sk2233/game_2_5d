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
	PlayerIdle = iota
	PlayerMove
	PlayerJump
	PlayerHurt
	PlayerLie
)

type Player0 struct {
	Pos      Vec
	Dir      float64
	Anim     *Animator
	Speed    float64
	State    int
	Actions  []func()
	LieTimer int
}

func (p *Player0) Update() {
	p.Anim.Update()
	p.Actions[p.State]()
}

func (p *Player0) move(speed Vec) {
	if speed.X != 0 {
		p.Dir = Sign(speed.X)
	}
	pos := p.Pos.Add(speed)
	pos.Z = RangeManager.GetZ(pos)
	if math.Abs(pos.Z-p.Pos.Z) <= Speed { // 在爬坡或者下坡
		p.Pos = pos
	} else if p.Pos.Z > pos.Z {
		p.Pos.X, p.Pos.Y = pos.X, pos.Y
	}
}

func (p *Player0) checkMove() bool {
	if GetAxis(ebiten.KeyD, ebiten.KeyA) != 0 || GetAxis(ebiten.KeyS, ebiten.KeyW) != 0 {
		p.State = PlayerMove
		p.Anim.Play("run")
		return true
	}
	return false
}

func (p *Player0) checkJump() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		p.State = PlayerJump
		p.Anim.RePlay("jump")
		p.Speed = 12
		return true
	}
	return false
}

func (p *Player0) checkIdle() bool {
	if GetAxis(ebiten.KeyD, ebiten.KeyA) == 0 && GetAxis(ebiten.KeyS, ebiten.KeyW) == 0 {
		p.State = PlayerIdle
		p.Anim.Play("idle")
		return true
	}
	return false
}

func (p *Player0) checkAir() bool {
	z := RangeManager.GetZ(p.Pos)
	if p.Pos.Z-z > Speed {
		p.State = PlayerJump
		p.Speed = 0
		p.Anim.RePlay("jump")
		p.Anim.SetFrame(1)
		return true
	}
	return false
}

func (p *Player0) Idle() {
	if BatchCheck(p.checkMove, p.checkJump) {
		return
	}
}

func (p *Player0) Move() {
	if BatchCheck(p.checkIdle, p.checkJump, p.checkAir) {
		return
	}
	p.move(Vec{X: GetAxis(ebiten.KeyD, ebiten.KeyA), Y: GetAxis(ebiten.KeyS, ebiten.KeyW)}.Scale(Speed))
}

func (p *Player0) Jump() {
	p.Speed -= G
	p.Pos.Z += p.Speed
	z := RangeManager.GetZ(p.Pos)
	if p.Pos.Z <= z {
		p.Pos.Z = z
		p.State = PlayerIdle
		p.Anim.RePlay("land")
		return
	}
	p.move(Vec{X: GetAxis(ebiten.KeyD, ebiten.KeyA), Y: GetAxis(ebiten.KeyS, ebiten.KeyW)}.Scale(Speed))
}

func (p *Player0) Hurt() {
	p.move(Vec{X: Speed * 2})
	p.Speed -= G
	p.Pos.Z += p.Speed
	z := RangeManager.GetZ(p.Pos)
	if p.Pos.Z <= z {
		p.Pos.Z = z
		p.State = PlayerLie
		p.Anim.RePlay("lie")
		p.LieTimer = 30
		return
	}
}

func (p *Player0) Lie() {
	if p.LieTimer > 0 {
		p.LieTimer--
	} else {
		p.State = PlayerIdle
		p.Anim.RePlay("idle")
	}
}

func (p *Player0) Draw(screen *ebiten.Image) {
	p.Anim.Draw(screen, ToVec2(p.Pos), p.Dir)
}

func (p *Player0) TriggerHurt() {
	p.Speed = 12
	p.State = PlayerHurt
	p.Anim.RePlay("hurt")
}

func (p *Player0) AnimEnd(anim *Anim) {
	switch anim.Name {
	case "land":
		p.Anim.Play("idle")
	}
}

func NewPlayer() *Player0 {
	res := &Player0{Anim: NewAnimator(), State: PlayerMove, Dir: 1}
	res.Anim.AddAnim(NewAnim("resources/player/hurt").SetPivot(PivotBottom))
	res.Anim.AddAnim(NewAnim("resources/player/idle").SetPivot(PivotBottom))
	res.Anim.AddAnim(NewAnim("resources/player/jump").SetPivot(PivotBottom).SetFps(10).FreezeFrame(1))
	res.Anim.AddAnim(NewAnim("resources/player/land").SetPivot(PivotBottom).SetFps(10))
	res.Anim.AddAnim(NewAnim("resources/player/lie").SetPivot(PivotBottom))
	res.Anim.AddAnim(NewAnim("resources/player/run").SetPivot(PivotBottom).SetFps(10))
	res.Anim.SetAnimEnd(res.AnimEnd)
	res.Pos = Vec{X: 50, Y: 100}
	res.Pos.Z = RangeManager.GetZ(res.Pos)
	res.Actions = make([]func(), 5)
	res.Actions[PlayerIdle] = res.Idle
	res.Actions[PlayerMove] = res.Move
	res.Actions[PlayerJump] = res.Jump
	res.Actions[PlayerHurt] = res.Hurt
	res.Actions[PlayerLie] = res.Lie
	return res
}
