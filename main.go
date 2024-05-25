package main

import (
	"game/game"
	"game/game/draw"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth*2, game.ScreenHeight*2)
	ebiten.SetWindowTitle("MyGame")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetScreenTransparent(false)

	gameState := &game.Game{}
	gameState.Squares = []*game.Square{}
	gameState.Squares = append(gameState.Squares,
		&game.Square{
			Size:      100,
			PosX:      200,
			PosY:      100,
			MovementX: 0.7,
			MovementY: 0.5,
			Colour:    draw.ColourBlue,
			Alpha:     90,
		}, &game.Square{
			Size:      100,
			PosX:      300,
			PosY:      200,
			MovementX: 0.4,
			MovementY: -0.9,
			Colour:    draw.ColourRed,
			Alpha:     90,
		}, &game.Square{
			Size:      100,
			PosX:      400,
			PosY:      200,
			MovementX: 0.3,
			MovementY: -0.5,
			Colour:    draw.ColourGreen,
			Alpha:     90,
		}, &game.Square{
			Size:      60,
			PosX:      330,
			PosY:      250,
			MovementX: 0.1,
			MovementY: -0.2,
			Colour:    draw.ColourPurple,
			Alpha:     170,
		},
	)
	gameState.Gravity = 0.05

	if err := ebiten.RunGame(gameState); err != nil {
		panic(err)
	}
}
