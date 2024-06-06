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
	CurrentPosX int
	CurrentPosY int

	Shown bool

	Image         *ebiten.Image
	Colour        color.Color
	ClickFunction func(int, int)
}

func (d *Divider) Init(vertical bool, position int, length int, colour color.Color) {
	xMin := 0
	xMax := length
	yMin := position
	yMax := position + 5
	if vertical {
		xMin = position
		xMax = position + 5
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
	d.CurrentPosX = xMin
	d.CurrentPosY = yMin
	d.Colour = colour
	d.Image = menuImage
}

func (d *Divider) UpdateSize(widthFactor, heightFactor int) {
	xMin := 0
	xMax := d.BaseLength * widthFactor
	yMin := d.BasePos * heightFactor
	yMax := d.BasePos*heightFactor + 5
	if d.Vertical {
		xMin = d.BasePos * widthFactor
		xMax = d.BasePos*widthFactor + 5
		yMin = 0
		yMax = d.BaseLength * heightFactor
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
