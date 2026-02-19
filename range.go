package main

import (
	"encoding/json"
	"os"

	"github.com/tidwall/gjson"
)

var (
	RangeManager = NewRangeManager()
)

type IRange interface {
	Update()
	Contains(pos Vec) bool
	GetZ(pos Vec) float64
}

type RangeManager0 struct {
	Ranges []IRange
}

func (m *RangeManager0) Update() {
	for _, item := range m.Ranges {
		item.Update()
	}
}

type Prop struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type Data struct {
	Type       string  `json:"type"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Width      float64 `json:"width"`
	Height     float64 `json:"height"`
	Properties []*Prop `json:"properties"`
}

func (m *RangeManager0) LoadRange(file string) {
	m.Ranges = make([]IRange, 0)
	bs, err := os.ReadFile(file)
	HandleErr(err)
	jsonStr := gjson.GetBytes(bs, "layers.1.objects").String()
	datas := make([]*Data, 0)
	err = json.Unmarshal([]byte(jsonStr), &datas)
	HandleErr(err)
	for _, data := range datas {
		left, right := parseHeight(data.Properties)
		switch data.Type {
		case "move":
			m.Ranges = append(m.Ranges, NewMoveRange(data.X, data.Y, data.Width, data.Height, left, right))
		case "trap":
			height := parseTrapHeight(data.Properties)
			m.Ranges = append(m.Ranges, NewTrapRange(data.X, data.Y, data.Width, data.Height, left, right, height))
		default:
			m.Ranges = append(m.Ranges, NewBaseRange(data.X, data.Y, data.Width, data.Height, left, right))
		}
	}
}

func (m *RangeManager0) GetZ(pos Vec) float64 {
	for _, item := range m.Ranges {
		if item.Contains(pos) {
			return item.GetZ(pos)
		}
	}
	return 10000
}

func parseHeight(props []*Prop) (float64, float64) {
	left, right := float64(0), float64(0)
	for _, prop := range props {
		switch prop.Name {
		case "left":
			left = prop.Value
		case "right":
			right = prop.Value
		case "height":
			return prop.Value, prop.Value
		}
	}
	return left, right
}

func parseTrapHeight(props []*Prop) float64 {
	for _, prop := range props {
		if prop.Name == "trap" {
			return prop.Value
		}
	}
	return 10000
}

func NewRangeManager() *RangeManager0 {
	res := &RangeManager0{}
	res.LoadRange("resources/map.tmj")
	return res
}

//===================BaseRange======================

type BaseRange struct {
	X, Y        float64
	W, H        float64
	Left, Right float64
}

func (b *BaseRange) Contains(pos Vec) bool {
	return pos.X >= b.X && pos.X <= b.X+b.W && pos.Y >= b.Y && pos.Y <= b.Y+b.H
}

func (b *BaseRange) GetZ(pos Vec) float64 {
	return (b.Right-b.Left)*(pos.X-b.X)/b.W + b.Left
}

func NewBaseRange(x float64, y float64, w float64, h float64, left float64, right float64) *BaseRange {
	return &BaseRange{X: x, Y: y, W: w, H: h, Left: left, Right: right}
}

func (b *BaseRange) Update() {
}

//=================TrapRange===================

type TrapRange struct {
	*BaseRange
	Height float64
}

func NewTrapRange(x float64, y float64, w float64, h float64, left float64, right float64, height float64) *TrapRange {
	return &TrapRange{BaseRange: NewBaseRange(x, y, w, h, left, right), Height: height}
}

func (t *TrapRange) Update() {
	if !t.Contains(Player.Pos) {
		return
	}
	z := t.GetZ(Player.Pos)
	if Player.Pos.Z >= z && Player.Pos.Z <= z+t.Height {
		Player.TriggerHurt()
	}
}

//================MoveRange==============

type MoveRange struct { // 默认认为高 8
	*BaseRange
}

func NewMoveRange(x float64, y float64, w float64, h float64, left float64, right float64) *MoveRange {
	return &MoveRange{BaseRange: NewBaseRange(x, y, w, h, left, right)}
}

func (m *MoveRange) Update() {
	if !m.Contains(Player.Pos) {
		return
	}
	z := m.GetZ(Player.Pos)
	if Player.Pos.Z >= z && Player.Pos.Z <= z+4 {
		Player.Pos.X--
		Player.Pos.Z = RangeManager.GetZ(Player.Pos)
	}
}
