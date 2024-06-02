package main

import (
	"game/game"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth*2, game.ScreenHeight*2)
	ebiten.SetWindowTitle("MyGame")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetScreenTransparent(false)

	if err := ebiten.RunGame(game.Init()); err != nil {
		panic(err)
	}
}
