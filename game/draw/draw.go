package draw

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Colour uint8

const (
	ColourWhite Colour = iota
	ColourRed
	ColourGreen
	ColourBlue
	ColourPurple
)

func getColourValues(colour Colour, alpha uint8) color.Color {
	switch colour {
	case ColourWhite:
		return color.NRGBA{
			R: 255,
			G: 255,
			B: 255,
			A: alpha,
		}
	case ColourRed:
		return color.NRGBA{
			R: 255,
			G: 0,
			B: 0,
			A: alpha,
		}
	case ColourGreen:
		return color.NRGBA{
			R: 0,
			G: 255,
			B: 0,
			A: alpha,
		}
	case ColourBlue:
		return color.NRGBA{
			R: 0,
			G: 0,
			B: 255,
			A: alpha,
		}
	case ColourPurple:
		return color.NRGBA{
			R: 50,
			G: 0,
			B: 255,
			A: alpha,
		}
	default:
		return color.NRGBA{
			R: 0,
			G: 0,
			B: 0,
			A: alpha,
		}
	}
}

func Square(screen *ebiten.Image, colour Colour, alpha uint8, size int, xPos, yPos float64) {
	square := ebiten.NewImage(size, size)
	square.Fill(getColourValues(colour, alpha))
	position := ebiten.GeoM{}
	position.Translate(xPos, yPos)

	screen.DrawImage(square, &ebiten.DrawImageOptions{
		GeoM:          position,
		CompositeMode: ebiten.CompositeModeSourceOver,
	})
}
