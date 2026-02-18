package main

import "math"

type Vec struct {
	X, Y, Z float64
}

func (v Vec) Add(t Vec) Vec {
	return Vec{X: v.X + t.X, Y: v.Y + t.Y, Z: v.Z + t.Z}
}

func (v Vec) Scale(s float64) Vec {
	return Vec{X: v.X * s, Y: v.Y * s, Z: v.Z * s}
}

func (v Vec) Snap(size float64) Vec {
	return Vec{X: math.Round(v.X/size) * size, Y: math.Round(v.Y/size) * size,
		Z: math.Round(v.Z/size) * size}
}
