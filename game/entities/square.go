package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
)

type Square struct {
	Size int

	PosX          float64
	MovementX     float64
	AccelerationX float64
	BounceX       float64

	PosY          float64
	MovementY     float64
	AccelerationY float64
	BounceY       float64

	Image *ebiten.Image
}

func NewSquare(red, green, blue, alpha uint8, size int, posX, posY, movX, movY, accX, accY, bounceX, bounceY float64) *Square {
	colour := color.NRGBA{
		R: red,
		G: green,
		B: blue,
		A: alpha,
	}
	square := ebiten.NewImage(size, size)
	square.Fill(colour)

	return &Square{
		Size:          size,
		PosX:          posX,
		MovementX:     movX,
		AccelerationX: accX,
		BounceX:       bounceX,
		PosY:          posY,
		MovementY:     movY,
		AccelerationY: accY,
		BounceY:       bounceY,
		Image:         square,
	}
}

func (s *Square) Click(x, y float64) {
	if s.PosX < x && s.PosX+float64(s.Size) > x && s.PosY < y && s.PosY+float64(s.Size) > y {
		s.AccelerationY = -math.Abs(s.AccelerationY)
		s.MovementY = -math.Abs(s.MovementY)
	}
}

func (s *Square) Update(screenWidth, screenHeight, timeDelta, gravity float64) {
	if s.PosX < 0 {
		s.MovementX *= -1
		s.AccelerationX *= -s.BounceX
		s.PosX = 0
	} else if s.PosX+float64(s.Size) > screenWidth {
		s.MovementX *= -1
		s.AccelerationX *= -s.BounceX
		s.PosX = screenWidth - float64(s.Size)
	}

	if s.PosY < 0 {
		s.MovementY *= -1
		s.AccelerationY *= -s.BounceY
		s.PosY = 0
	} else if s.PosY+float64(s.Size) > screenHeight {
		s.MovementY *= -1
		s.AccelerationY *= -s.BounceY
		s.PosY = screenHeight - float64(s.Size)
	}

	s.AccelerationY += gravity * timeDelta
	s.PosX += (s.MovementX + s.AccelerationX) * timeDelta
	s.PosY += (s.MovementY + s.AccelerationY) * timeDelta
}
