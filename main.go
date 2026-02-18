package main

import "github.com/hajimehoshi/ebiten/v2"

func main() {
	ebiten.SetWindowSize(Width*4, Height*4)
	ebiten.SetWindowTitle("Nekketsu Fighting Legend")
	err := ebiten.RunGame(NewGame())
	HandleErr(err)
}
