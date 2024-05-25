package game

import (
	"game/game/draw"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

type Game struct {
	touchIDs []ebiten.TouchID
	op       ebiten.DrawImageOptions

	Squares []*Square
	Gravity float64
}

type Square struct {
	Size          int
	PosX          float64
	PosY          float64
	MovementX     float64
	AccelerationX float64
	MovementY     float64
	AccelerationY float64
	Colour        draw.Colour
	Alpha         uint8
}

func (g *Game) Update() error {
	g.touchIDs = ebiten.AppendTouchIDs(g.touchIDs[:0])

	for _, square := range g.Squares {
		if square.PosX < 0 || square.PosX+float64(square.Size) > ScreenWidth {
			square.MovementX *= -1
			square.AccelerationX *= -1
		}
		if square.PosY < 0 || square.PosY+float64(square.Size) > ScreenHeight {
			square.MovementY *= -1
			square.AccelerationY *= -1
		}

		square.AccelerationY += g.Gravity
		square.PosX += square.MovementX + square.AccelerationX
		square.PosY += square.MovementY + square.AccelerationY
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Test debug")

	for _, square := range g.Squares {
		draw.Square(screen, square.Colour, square.Alpha, square.Size, square.PosX, square.PosY)
	}
}
