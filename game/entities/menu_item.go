package entities

import (
	"game/game/helper"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"image/color"
)

type MenuItem struct {
	BaseWidth    float64
	BaseHeight   float64
	BasePosX     float64
	BasePosY     float64
	BaseTextSize float64

	CurrentWidth  float64
	CurrentHeight float64
	CurrentPosX   float64
	CurrentPosY   float64

	Shown bool

	Image            *ebiten.Image
	BackgroundColour color.Color
	Colour           color.Color
	TextColour       color.Color
	Name             string
	Font             *text.GoTextFaceSource
	ClickFunction    func(int, int)
}

func (mi *MenuItem) Init(width, height, x, y float64, shown bool, name string, font *text.GoTextFaceSource, textSize float64, textColour, backgroundColour, colour color.Color, clickFunction func(int, int)) {
	menuImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: int(x), Y: int(y)},
		Max: image.Point{X: int(x + width), Y: int(y + height)},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	menuImage.Fill(backgroundColour)
	for i := 0; i <= int(width+x); i++ {
		for n := 0; n <= 5; n++ {
			menuImage.Set(i, int(y)+n, colour)
			menuImage.Set(i, int(y+height)-n, colour)
		}
	}
	for i := 0; i <= int(height+y); i++ {
		for n := 0; n <= 5; n++ {
			menuImage.Set(int(x)+n, i, colour)
			menuImage.Set(int(x+width)-n, i, colour)
		}
	}

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(textColour)
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter

	middleX := width/2 + x
	middleY := height/2 + y
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

func (mi *MenuItem) UpdateSize(widthFactor, heightFactor float64) {
	newWidth, newHeight, newX, newY := helper.GetNewSizeAndPosition(mi.BaseWidth, mi.BaseHeight, mi.BasePosX, mi.BasePosY, widthFactor, heightFactor, 0, 0)
	newTextSize := helper.GetNewTextSize(mi.BaseTextSize, heightFactor, newWidth, mi.Name)

	menuImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: int(newX), Y: int(newY)},
		Max: image.Point{X: int(newX + newWidth), Y: int(newY + newHeight)},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	menuImage.Fill(mi.BackgroundColour)
	for i := 0; i <= int(newWidth+newX); i++ {
		for n := 0; n <= 5; n++ {
			menuImage.Set(i, int(newY)+n, mi.Colour)
			menuImage.Set(i, int(newY+newHeight)-n, mi.Colour)
		}
	}
	for i := 0; i <= int(newHeight+newY); i++ {
		for n := 0; n <= 5; n++ {
			menuImage.Set(int(newX)+n, i, mi.Colour)
			menuImage.Set(int(newX+newWidth)-n, i, mi.Colour)
		}
	}

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(mi.TextColour)
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter

	middleX := newWidth/2 + newX
	middleY := newHeight/2 + newY
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

	if mi.CurrentPosX < float64(x) && mi.CurrentPosX+mi.CurrentWidth > float64(x) &&
		mi.CurrentPosY < float64(y) && mi.CurrentPosY+mi.CurrentHeight > float64(y) {
		mi.ClickFunction(x, y)
	}
}
