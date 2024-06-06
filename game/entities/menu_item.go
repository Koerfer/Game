package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"image/color"
	"math"
)

type MenuItem struct {
	BaseWidth    int
	BaseHeight   int
	BasePosX     int
	BasePosY     int
	BaseTextSize float64

	CurrentWidth  int
	CurrentHeight int
	CurrentPosX   int
	CurrentPosY   int

	Shown bool

	Image            *ebiten.Image
	BackgroundColour color.Color
	Colour           color.Color
	TextColour       color.Color
	Name             string
	Font             *text.GoTextFaceSource
	ClickFunction    func(int, int)
}

func (mi *MenuItem) Init(width, height, x, y int, shown bool, name string, font *text.GoTextFaceSource, textSize float64, textColour, backgroundColour, colour color.Color, clickFunction func(int, int)) {
	menuImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: x, Y: y},
		Max: image.Point{X: x + width, Y: y + height},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	menuImage.Fill(backgroundColour)
	for i := 0; i <= width+x; i++ {
		for n := 0; n <= 5; n++ {
			menuImage.Set(i, y+n, colour)
			menuImage.Set(i, y+height-n, colour)
		}
	}
	for i := 0; i <= height+y; i++ {
		for n := 0; n <= 5; n++ {
			menuImage.Set(x+n, i, colour)
			menuImage.Set(x+width-n, i, colour)
		}
	}

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(textColour)
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter

	middleX := float64(width)/2 + float64(x)
	middleY := float64(height)/2 + float64(y)
	op.GeoM.Translate(middleX, middleY)

	text.Draw(menuImage, name, &text.GoTextFace{
		Source: font,
		Size:   textSize,
	}, op)

	mi.BaseWidth = width
	mi.BaseHeight = height
	mi.BasePosX = x
	mi.BasePosY = y
	mi.CurrentWidth = width
	mi.CurrentHeight = height
	mi.CurrentPosX = x
	mi.CurrentPosY = y
	mi.Shown = shown
	mi.Image = menuImage
	mi.BackgroundColour = backgroundColour
	mi.Colour = colour
	mi.TextColour = textColour
	mi.Name = name
	mi.Font = font
	mi.BaseTextSize = textSize
	mi.ClickFunction = clickFunction

}

func (mi *MenuItem) UpdateSize(widthFactor, heightFactor int) {
	newWidth := mi.BaseWidth * widthFactor
	newHeight := mi.BaseHeight * heightFactor
	newX := mi.BasePosX * widthFactor
	newY := mi.BasePosY * heightFactor
	newTextSize := mi.BaseTextSize * float64(heightFactor)
	if newTextSize*4 > float64(newWidth) {
		newTextSize = math.Min(newTextSize, float64(newWidth)/4)
	}

	menuImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: newX, Y: newY},
		Max: image.Point{X: newX + newWidth, Y: newY + newHeight},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	menuImage.Fill(mi.BackgroundColour)
	for i := 0; i <= newWidth+newX; i++ {
		for n := 0; n <= 5; n++ {
			menuImage.Set(i, newY+n, mi.Colour)
			menuImage.Set(i, newY+newHeight-n, mi.Colour)
		}
	}
	for i := 0; i <= newHeight+newY; i++ {
		for n := 0; n <= 5; n++ {
			menuImage.Set(newX+n, i, mi.Colour)
			menuImage.Set(newX+newWidth-n, i, mi.Colour)
		}
	}

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(mi.TextColour)
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter

	middleX := float64(newWidth)/2 + float64(newX)
	middleY := float64(newHeight)/2 + float64(newY)
	op.GeoM.Translate(middleX, middleY)

	text.Draw(menuImage, mi.Name, &text.GoTextFace{
		Source: mi.Font,
		Size:   newTextSize,
	}, op)

	mi.CurrentWidth = newWidth
	mi.CurrentHeight = newHeight
	mi.CurrentPosX = newX
	mi.CurrentPosY = newY
	mi.Image = menuImage
}

func (mi *MenuItem) Click(x, y int) {
	if !mi.Shown {
		return
	}

	if mi.CurrentPosX < x && mi.CurrentPosX+mi.CurrentWidth > x &&
		mi.CurrentPosY < y && mi.CurrentPosY+mi.CurrentHeight > y {
		mi.ClickFunction(x, y)
	}
}
