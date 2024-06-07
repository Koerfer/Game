package cards

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"image/color"
	"math"
)

type Card struct {
	Name                    string
	Description             string
	BaseNameTextSize        float64
	BaseDescriptionTextSize float64
	TextColour              color.Color
	Font                    *text.GoTextFaceSource

	BaseWidth  int
	BaseHeight int
	BasePosX   int
	BasePosY   int

	CurrentWidth  int
	CurrentHeight int
	CurrentPosX   int
	CurrentPosY   int

	BackgroundColour color.Color
	Colour           color.Color
	Image            *ebiten.Image
}

func (c *Card) Init(width, height, x, y int, name, description string, font *text.GoTextFaceSource, nameTextSize, descriptionTextSize float64, textColour, backgroundColour, colour color.Color) {
	cardImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: x, Y: y},
		Max: image.Point{X: x + width, Y: y + height},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	cardImage.Fill(backgroundColour)
	for i := 0; i <= width+x; i++ {
		for n := 0; n <= 5; n++ {
			cardImage.Set(i, y+n, colour)
			cardImage.Set(i, y+height-n, colour)
		}
	}
	for i := 0; i <= height+y; i++ {
		for n := 0; n <= 5; n++ {
			cardImage.Set(x+n, i, colour)
			cardImage.Set(x+width-n, i, colour)
		}
	}

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(textColour)
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter

	middleX := float64(width)/2 + float64(x)
	middleY := float64(height)/2 + float64(y)
	op.GeoM.Translate(middleX, middleY)

	text.Draw(cardImage, name, &text.GoTextFace{
		Source: font,
		Size:   nameTextSize,
	}, op)

	c.BaseWidth = width
	c.BaseHeight = height
	c.BasePosX = x
	c.BasePosY = y
	c.CurrentWidth = width
	c.CurrentHeight = height
	c.CurrentPosX = x
	c.CurrentPosY = y
	c.Image = cardImage
	c.BackgroundColour = backgroundColour
	c.Colour = colour
	c.TextColour = textColour
	c.Name = name
	c.Font = font
	c.BaseNameTextSize = nameTextSize
	c.BaseDescriptionTextSize = descriptionTextSize
}

func (c *Card) UpdateSize(widthFactor, heightFactor float64) {
	shift := 12
	if c.BasePosX == 8 {
		shift = 6
	} else if c.BasePosX == 12 {
		shift = 0
	}
	newWidth := int(float64(c.BaseWidth) * widthFactor)
	newHeight := int(float64(c.BaseHeight) * heightFactor)
	newX := int(float64(c.BasePosX)*widthFactor) + shift
	newY := int(float64(c.BasePosY) * heightFactor)
	newNameTextSize := c.BaseNameTextSize * heightFactor
	//newDescriptionTextSize := c.BaseDescriptionTextSize * float64(heightFactor)
	if newNameTextSize*float64(len(c.Name))*3.6/5 > float64(newWidth) {
		newNameTextSize = math.Min(newNameTextSize, 1.4*float64(newWidth)/float64(len(c.Name)))
	}

	cardImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: newX, Y: newY},
		Max: image.Point{X: newX + newWidth, Y: newY + newHeight},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	cardImage.Fill(c.BackgroundColour)
	for i := 0; i <= newWidth+newX; i++ {
		for n := 0; n <= 5; n++ {
			cardImage.Set(i, newY+n, c.Colour)
			cardImage.Set(i, newY+newHeight-n, c.Colour)
		}
	}
	for i := 0; i <= newHeight+newY; i++ {
		for n := 0; n <= 5; n++ {
			cardImage.Set(newX+n, i, c.Colour)
			cardImage.Set(newX+newWidth-n, i, c.Colour)
		}
	}

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(c.TextColour)
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter

	middleX := float64(newWidth)/2 + float64(newX)
	middleY := float64(newHeight)/2 + float64(newY)
	op.GeoM.Translate(middleX, middleY)

	text.Draw(cardImage, c.Name, &text.GoTextFace{
		Source: c.Font,
		Size:   newNameTextSize,
	}, op)

	c.CurrentWidth = newWidth
	c.CurrentHeight = newHeight
	c.CurrentPosX = newX
	c.CurrentPosY = newY
	c.Image = cardImage
}
