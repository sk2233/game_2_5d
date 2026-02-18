package main

var (
	Camera = NewCamera0()
)

type Camera0 struct {
	X float64 // 只控制 x 即可
}

func (c *Camera0) Update() {
	c.X = Player.Pos.X - Width/2
	if c.X < 0 {
		c.X = 0
	}
	if c.X > 496-Width {
		c.X = 496 - Width
	}
}

func NewCamera0() *Camera0 {
	return &Camera0{}
}
