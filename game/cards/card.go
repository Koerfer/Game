package cards

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"image/color"
	"time"
)

type Card struct {
	Name                    string
	Description             string
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

	ActiveSingleTargetDamageBoost float64
	PassiveDamageBoost            float64
	ActivationTime                time.Duration
	CoolDown                      time.Duration
	PlayCard                      *PlayCard

	State State
}

type State uint8

const (
	StateLocked State = iota
	StateUnlocked
	StateSelected
)

func (c *Card) Init(x, y float64, name string, font *text.GoTextFaceSource, nameTextSize, descriptionTextSize float64, textColour, backgroundColour, colour color.Color) {
	whiteImage.Fill(color.White)

	cardImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: int(x), Y: int(y)},
		Max: image.Point{X: int(x + 4), Y: int(y + 8)},
	}, &ebiten.NewImageOptions{Unmanaged: false})

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

	c.State = StateUnlocked

	c.addPlayCard()
}

func (c *Card) addEffect(description string, damageBoost, singleTargetDamageBoost float64, activationTime, coolDown time.Duration) {
	if damageBoost != 1 {
		c.Description = fmt.Sprintf(description, int(damageBoost), int(singleTargetDamageBoost), int(coolDown.Seconds()))
	} else {
		c.Description = description
	}

	c.PassiveDamageBoost = damageBoost
	c.ActiveSingleTargetDamageBoost = singleTargetDamageBoost
	c.ActivationTime = activationTime
	c.CoolDown = coolDown
	c.PlayCard.CoolDownRemaining = 0
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
