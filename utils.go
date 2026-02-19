package main

import (
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	imageCache = make(map[string]*ebiten.Image)
)

func LoadImage(name string) *ebiten.Image {
	name = path.Clean(name)
	if _, ok := imageCache[name]; !ok {
		img, _, err := ebitenutil.NewImageFromFile(name)
		HandleErr(err)
		imageCache[name] = img
	}
	return imageCache[name]
}

func LoadImages(name string) []*ebiten.Image {
	items, err := os.ReadDir(name)
	HandleErr(err)
	res := make([]*ebiten.Image, 0)
	for _, item := range items {
		res = append(res, LoadImage(path.Join(name, item.Name())))
	}
	return res
}

func ToVec2(vec Vec) Vec {
	return Vec{X: vec.X, Y: Height - 112/2 + vec.Y/2 - vec.Z/2}
}

func GetAxis(pos, neg ebiten.Key) float64 {
	if ebiten.IsKeyPressed(pos) {
		return 1
	}
	if ebiten.IsKeyPressed(neg) {
		return -1
	}
	return 0
}

func Sign(val float64) float64 {
	if val > 0 {
		return 1
	} else if val < 0 {
		return -1
	} else {
		return 0
	}
}

func BatchCheck(checkers ...func() bool) bool {
	for _, checker := range checkers {
		if checker() {
			return true
		}
	}
	return false
}
