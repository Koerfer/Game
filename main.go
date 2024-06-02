package main

import (
	"game/game"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
)

func main() {
	ebiten.SetWindowSize(900, 900)
	ebiten.SetWindowSizeLimits(900, 900, -1, -1)
	ebiten.SetWindowTitle("MyGame")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game.Init()); err != nil {
		panic(err)
	}
}
