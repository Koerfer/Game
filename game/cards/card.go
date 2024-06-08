package cards

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"image/color"
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
	SelectedImage    *ebiten.Image
	UnlockedImage    *ebiten.Image
	LockedImage      *ebiten.Image

	State State
}

type State uint8

const (
	StateLocked State = iota
	StateUnlocked
	StateSelected
)

func (c *Card) Init(width, height, x, y float64, name, description string, font *text.GoTextFaceSource, nameTextSize, descriptionTextSize float64, textColour, backgroundColour, colour color.Color) {
	whiteImage.Fill(color.White)

	cardImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: int(x), Y: int(y)},
		Max: image.Point{X: int(x + width), Y: int(y + height)},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	c.BaseWidth = width
	c.BaseHeight = height
	c.BasePosX = x
	c.BasePosY = y
	c.CurrentWidth = width
	c.CurrentHeight = height
	c.CurrentPosX = x
	c.CurrentPosY = y
	c.SelectedImage = cardImage
	c.LockedImage = cardImage
	c.UnlockedImage = cardImage
	c.BackgroundColour = backgroundColour
	c.Colour = colour
	c.TextColour = textColour
	c.Name = name
	c.Description = description
	c.Font = font
	c.BaseNameTextSize = nameTextSize
	c.BaseDescriptionTextSize = descriptionTextSize

	c.State = StateUnlocked
}

func (c *Card) Click(x, y, numberSelected int) int {
	if c.CurrentPosX < float64(x) && c.CurrentPosX+c.CurrentWidth > float64(x) &&
		c.CurrentPosY < float64(y) && c.CurrentPosY+c.CurrentHeight > float64(y) {
		switch c.State {
		case StateLocked:
			return 0
		case StateSelected:
			c.State = StateUnlocked
			return -1
		case StateUnlocked:
			if numberSelected < 3 {
				c.State = StateSelected
				return 1
			}
			return 0
		}
	}

	return 0
}
