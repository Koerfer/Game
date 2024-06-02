package main

import (
	"game/game"
	"game/game/entities"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
	"time"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth*2, game.ScreenHeight*2)
	ebiten.SetWindowTitle("MyGame")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetScreenTransparent(false)

	gameState := &game.Game{}
	gameState.Squares = []*entities.Square{}
	gameState.Squares = append(gameState.Squares,
		entities.NewSquare(255, 0, 0, 90, 100, 300, 200, 0.11, -0.09, 0, 0, 1, 0.9),
		entities.NewSquare(0, 255, 0, 90, 120, 100, 500, 0.05, 0.1, 0, 0, 1, 0.9),
	)
	gameState.Gravity = 0.0005
	gameState.PreviousUpdateTime = time.Now()

	if err := ebiten.RunGame(gameState); err != nil {
		panic(err)
	}
}
