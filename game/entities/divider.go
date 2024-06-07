package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

type Divider struct {
	Vertical bool

	BasePos     int
	BaseLength  int
	Width       int
	Align       Align
	CurrentPosX int
	CurrentPosY int

	Shown bool

	Image         *ebiten.Image
	Colour        color.Color
	ClickFunction func(int, int)
}

type Align uint8

const (
	Middle Align = iota
	Left
	Right
)

func (d *Divider) Init(vertical bool, position, length, width int, align Align, colour color.Color) {
	var leftUp int
	var rightDown int

	switch align {
	case Middle:
		leftUp = width / 2
		rightDown = width / 2
	case Left:
		leftUp = 0
		rightDown = width
	case Right:
		leftUp = width
		rightDown = 0
	}

	xMin := 0
	xMax := length
	yMin := position - leftUp
	yMax := position + rightDown
	if vertical {
		xMin = position - leftUp
		xMax = position + rightDown
		yMin = 0
		yMax = length
	}

	menuImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: xMin, Y: yMin},
		Max: image.Point{X: xMax, Y: yMax},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	menuImage.Fill(colour)

	d.Vertical = vertical
	d.BasePos = position
	d.BaseLength = length
	d.Width = width
	d.Align = align
	d.CurrentPosX = xMin
	d.CurrentPosY = yMin
	d.Colour = colour
	d.Image = menuImage
}

func (d *Divider) UpdateSize(widthFactor, heightFactor float64) {
	var leftUp float64
	var rightDown float64

	switch d.Align {
	case Middle:
		leftUp = float64(d.Width) / 2
		rightDown = float64(d.Width) / 2
	case Left:
		leftUp = 0
		rightDown = float64(d.Width)
	case Right:
		leftUp = float64(d.Width)
		rightDown = 0
	}

	xMin := 0
	xMax := int(float64(d.BaseLength) * widthFactor)
	yMin := int(float64(d.BasePos)*heightFactor - leftUp)
	yMax := int(float64(d.BasePos)*heightFactor + rightDown)
	if d.Vertical {
		xMin = int(float64(d.BasePos)*widthFactor - leftUp)
		xMax = int(float64(d.BasePos)*widthFactor + rightDown)
		yMin = 0
		yMax = int(float64(d.BaseLength) * heightFactor)
	}

	menuImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: xMin, Y: yMin},
		Max: image.Point{X: xMax, Y: yMax},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	menuImage.Fill(d.Colour)

	d.CurrentPosX = xMin
	d.CurrentPosY = yMin
	d.Image = menuImage
}
