package main

import (
	"game/game"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
)

func main() {
	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowSizeLimits(900, 900, -1, -1)
	ebiten.SetWindowTitle("MyGame")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game.Init()); err != nil {
		panic(err)
	}
}
