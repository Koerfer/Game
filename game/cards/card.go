package cards

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"image/color"
)

type Card struct {
	Id int

	Name                    string
	BaseNameTextSize        float64
	BaseDescriptionTextSize float64
	TextColour              color.Color
	Font                    *text.GoTextFaceSource

	baseWidth  float64
	baseHeight float64
	basePosX   float64
	basePosY   float64

	currentWidth  float64
	currentHeight float64
	CurrentPosX   float64
	CurrentPosY   float64

	BackgroundColour color.Color
	Colour           color.Color
	SelectedImage    *ebiten.Image
	UnlockedImage    *ebiten.Image
	LockedImage      *ebiten.Image

	PlayCard *PlayCard

	State State
}

type State uint8

const (
	StateLocked State = iota
	StateUnlocked
	StateSelected
)

func (c *Card) Init(id int, x, y float64, name string, font *text.GoTextFaceSource, nameTextSize, descriptionTextSize float64,
	textColour, backgroundColour, colour color.Color) {
	whiteImage.Fill(color.White)

	cardImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: int(x), Y: int(y)},
		Max: image.Point{X: int(x + 4), Y: int(y + 8)},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	c.Id = id

	c.baseWidth = 4
	c.baseHeight = 8
	c.basePosX = x
	c.basePosY = y
	c.currentWidth = 4
	c.currentHeight = 8
	c.CurrentPosX = x
	c.CurrentPosY = y
	c.SelectedImage = cardImage
	c.LockedImage = cardImage
	c.UnlockedImage = cardImage
	c.BackgroundColour = backgroundColour
	c.Colour = colour
	c.TextColour = textColour
	c.Name = name
	c.Font = font
	c.BaseNameTextSize = nameTextSize
	c.BaseDescriptionTextSize = descriptionTextSize

	c.State = StateLocked

	c.addPlayCard()
}

func (c *Card) Click(x, y, numberSelected int) int {
	if c.CurrentPosX < float64(x) && c.CurrentPosX+c.currentWidth > float64(x) &&
		c.CurrentPosY < float64(y) && c.CurrentPosY+c.currentHeight > float64(y) {
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

func (c *Card) Upgrade(x, y, upgrades int) bool {
	if c.CurrentPosX < float64(x) && c.CurrentPosX+c.currentWidth > float64(x) &&
		c.CurrentPosY < float64(y) && c.CurrentPosY+c.currentHeight > float64(y) {
		switch c.State {
		case StateLocked:
			return false
		default:
			if upgrades == 0 {
				return false
			}
			if c.PlayCard.ActiveMultiTargetBoost != 1 {
				c.PlayCard.ActiveMultiTargetBoost *= 2
			}
			if c.PlayCard.PassiveMultiTargetBoost != 1 {
				c.PlayCard.PassiveMultiTargetBoost *= 2
			}
			if c.PlayCard.ActiveSingleTargetDamageBoost != 1 {
				c.PlayCard.ActiveSingleTargetDamageBoost *= 2
			}
			if c.PlayCard.PassiveDamageBoost != 1 {
				c.PlayCard.PassiveDamageBoost *= 2
			}

			return true
		}
	}

	return false
}
