package cards

import (
	"game/game/helper"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"image/color"
	"strings"
)

type Card struct {
	Name                    string
	Description             string
	BaseNameTextSize        float64
	BaseDescriptionTextSize float64
	TextColour              color.Color
	Font                    *text.GoTextFaceSource

	BaseWidth  float64
	BaseHeight float64
	BasePosX   float64
	BasePosY   float64

	CurrentWidth  float64
	CurrentHeight float64
	CurrentPosX   float64
	CurrentPosY   float64

	BackgroundColour color.Color
	Colour           color.Color
	Image            *ebiten.Image
}

func (c *Card) Init(width, height, x, y float64, name, description string, font *text.GoTextFaceSource, nameTextSize, descriptionTextSize float64, textColour, backgroundColour, colour color.Color) {
	cardImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: int(x), Y: int(y)},
		Max: image.Point{X: int(x + width), Y: int(y + height)},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	cardImage.Fill(backgroundColour)
	for i := 0; i <= int(width+x); i++ {
		for n := 0; n <= 5; n++ {
			cardImage.Set(i, int(y)+n, colour)
			cardImage.Set(i, int(y+height)-n, colour)
		}
	}
	for i := 0; i <= int(height+y); i++ {
		for n := 0; n <= 5; n++ {
			cardImage.Set(int(x)+n, i, colour)
			cardImage.Set(int(x+width)-n, i, colour)
		}
	}

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(textColour)
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter

	middleX := width/2 + x
	middleY := height/2 + y
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
	c.Description = description
	c.Font = font
	c.BaseNameTextSize = nameTextSize
	c.BaseDescriptionTextSize = descriptionTextSize
}

func (c *Card) UpdateSize(widthFactor, heightFactor float64) {
	shiftX := 12.0
	if c.BasePosX == 8 {
		shiftX = 5.0
	} else if c.BasePosX == 12 {
		shiftX = -2.0
	}

	shiftY := 3.0
	if c.BasePosY == 8 {
		shiftY = -2.0
	}

	newWidth, newHeight, newX, newY := helper.GetNewSizeAndPosition(c.BaseWidth, c.BaseHeight, c.BasePosX, c.BasePosY, widthFactor, heightFactor, shiftX, shiftY)
	newNameTextSize := helper.GetNewTextSize(c.BaseNameTextSize, heightFactor, newWidth, c.Name)

	cardImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: int(newX), Y: int(newY)},
		Max: image.Point{X: int(newX + newWidth), Y: int(newY + newHeight)},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	cardImage.Fill(c.BackgroundColour)
	for i := 0; i <= int(newWidth+newX); i++ {
		for n := 0; n <= 5; n++ {
			cardImage.Set(i, int(newY)+n, c.Colour)
			cardImage.Set(i, int(newY+newHeight)-n, c.Colour)
		}
	}
	for i := 0; i <= int(newHeight+newY); i++ {
		for n := 0; n <= 5; n++ {
			cardImage.Set(int(newX)+n, i, c.Colour)
			cardImage.Set(int(newX+newWidth)-n, i, c.Colour)
		}
	}

	c.printName(cardImage, newNameTextSize, newWidth, newX, newY)
	c.printDescription(cardImage, newNameTextSize, heightFactor, newWidth, newX, newY)

	c.CurrentWidth = newWidth
	c.CurrentHeight = newHeight
	c.CurrentPosX = newX
	c.CurrentPosY = newY
	c.Image = cardImage
}

func (c *Card) printName(cardImage *ebiten.Image, newNameTextSize, newWidth, newX, newY float64) {
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(c.TextColour)
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter

	posX := newWidth/2 + newX
	posY := newY + newNameTextSize/1.85
	op.GeoM.Translate(posX, posY)

	text.Draw(cardImage, c.Name, &text.GoTextFace{
		Source: c.Font,
		Size:   newNameTextSize,
	}, op)
}

func (c *Card) printDescription(cardImage *ebiten.Image, newNameTextSize, heightFactor, newWidth, newX, newY float64) {
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(c.TextColour)
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter

	posX := newWidth/2 + newX
	posY := newY + newNameTextSize*1.4
	op.GeoM.Translate(posX, posY)

	newDescriptionTextSize := helper.GetNewTextSize(c.BaseDescriptionTextSize, heightFactor, newWidth, c.Description)
	widthOfText, _ := text.Measure(c.Description, &text.GoTextFace{
		Source: c.Font,
		Size:   newDescriptionTextSize,
	}, 0)

	splitDescription := strings.Split(c.Description, "\n")
	for _, subDescription := range splitDescription {
		text.Draw(cardImage, subDescription, &text.GoTextFace{
			Source: c.Font,
			Size:   newWidth / widthOfText * newDescriptionTextSize * 0.9,
		}, op)
		op.GeoM.Translate(0, newWidth/widthOfText*newDescriptionTextSize)
	}
}
