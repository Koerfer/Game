package cards

import (
	"fmt"
	"game/game/colours"
	"game/game/helper"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"strings"
	"time"
)

type PlayCard struct {
	Id int

	baseWidth  float64
	baseHeight float64
	basePosX   float64
	basePosY   float64

	CurrentWidth  float64
	currentHeight float64
	CurrentPosX   float64
	CurrentPosY   float64

	Description       string
	BaseDescription   string
	Active            bool
	ActiveTime        time.Duration
	ActiveRemaining   time.Duration
	CoolDown          time.Duration
	CoolDownRemaining time.Duration

	ActiveSingleTargetDamageBoost float64
	PassiveDamageBoost            float64

	ActiveMultiTargetBoost  int
	PassiveMultiTargetBoost int

	PassiveTimeSlow int64
	ActiveTimeSkip  time.Duration

	PlayImage *ebiten.Image
}

func (c *Card) addPlayCard() {
	playCardImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 4, Y: 5},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	c.PlayCard = &PlayCard{
		Id:                c.Id,
		baseWidth:         4,
		baseHeight:        5,
		basePosX:          0,
		basePosY:          0,
		Active:            false,
		CoolDownRemaining: 0,
		PlayImage:         playCardImage,
	}
}

func (pc *PlayCard) addEffect(description string, passiveDamageBoost, activeSingleTargetDamageBoost float64,
	passiveMultiBoost, activeMultiBoost int,
	passiveTimeSlow int64, activeTimeSkip time.Duration, activationTime, coolDown time.Duration) {
	pc.BaseDescription = description

	if passiveDamageBoost != 1 {
		pc.Description = fmt.Sprintf(description, int(passiveDamageBoost), int(activeSingleTargetDamageBoost), int(coolDown.Seconds()))
	} else if passiveMultiBoost != 1 {
		pc.Description = fmt.Sprintf(description, passiveMultiBoost, activeMultiBoost, int(coolDown.Seconds()))
	} else if passiveTimeSlow != 1 {
		pc.Description = fmt.Sprintf(description, int(passiveTimeSlow), int(activeTimeSkip.Seconds()), int(coolDown.Seconds()))
	} else {
		pc.Description = description
	}

	pc.PassiveDamageBoost = passiveDamageBoost
	pc.ActiveSingleTargetDamageBoost = activeSingleTargetDamageBoost

	pc.PassiveMultiTargetBoost = passiveMultiBoost
	pc.ActiveMultiTargetBoost = activeMultiBoost

	pc.PassiveTimeSlow = passiveTimeSlow
	pc.ActiveTimeSkip = activeTimeSkip

	pc.ActiveTime = activationTime
	pc.CoolDown = coolDown
	pc.CoolDownRemaining = 0
}

func (c *Card) updatePlayCard(widthFactor, heightFactor float64) {
	if c.PlayCard.PassiveDamageBoost != 1 {
		c.PlayCard.Description = strings.Replace(fmt.Sprintf(c.PlayCard.BaseDescription, c.PlayCard.PassiveDamageBoost, c.PlayCard.ActiveSingleTargetDamageBoost, int(c.PlayCard.CoolDown.Seconds())), "e+0", "e", 2)
		c.PlayCard.Description = strings.Replace(c.PlayCard.Description, "e+", "e", 2)
	} else if c.PlayCard.PassiveMultiTargetBoost != 1 {
		c.PlayCard.Description = strings.Replace(fmt.Sprintf(c.PlayCard.BaseDescription, float64(c.PlayCard.PassiveMultiTargetBoost), c.PlayCard.ActiveMultiTargetBoost, int(c.PlayCard.CoolDown.Seconds())), "e+0", "e", 1)
		c.PlayCard.Description = strings.Replace(c.PlayCard.Description, "e+", "e", 1)
	} else if c.PlayCard.PassiveTimeSlow != 1 {
		c.PlayCard.Description = fmt.Sprintf(c.PlayCard.BaseDescription, c.PlayCard.PassiveTimeSlow, int(c.PlayCard.ActiveTimeSkip.Seconds()), int(c.PlayCard.CoolDown.Seconds()))
	}

	shiftX := 12.0
	if c.PlayCard.basePosX == 8 {
		shiftX = 8.0
	} else if c.PlayCard.basePosX == 12 {
		shiftX = 4.0
	}

	newWidth, newHeight, newX, newY := helper.GetNewSizeAndPosition(c.PlayCard.baseWidth, c.PlayCard.baseHeight, c.PlayCard.basePosX, c.PlayCard.basePosY, widthFactor, heightFactor, shiftX, 0)
	newWidth -= 8
	newHeight -= 7

	newNameTextSize := helper.GetNewTextSize(c.BaseNameTextSize, heightFactor, newWidth, c.Name)

	backgroundColour := colours.DarkBlue
	colour := c.Colour
	textColour := c.TextColour

	playCardImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: int(newX), Y: int(newY)},
		Max: image.Point{X: int(newX + newWidth), Y: int(newY + newHeight)},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	playCardImage.Fill(backgroundColour)
	for i := 0; i <= int(newWidth+newX); i++ {
		for n := 0; n <= 5; n++ {
			playCardImage.Set(i, int(newY)+n, colour)
			playCardImage.Set(i, int(newY+newHeight)-n, colour)
		}
	}
	for i := 0; i <= int(newHeight+newY); i++ {
		for n := 0; n <= 5; n++ {
			playCardImage.Set(int(newX)+n, i, colour)
			playCardImage.Set(int(newX+newWidth)-n, i, colour)
		}
	}

	c.PlayCard.CurrentWidth = newWidth
	c.PlayCard.currentHeight = newHeight
	c.PlayCard.CurrentPosX = newX
	c.PlayCard.CurrentPosY = newY
	c.printNamePlay(playCardImage, newNameTextSize, newWidth, newHeight, newX, newY, textColour)
	c.printDescription(playCardImage, heightFactor, newWidth, newHeight, newX, newY, textColour)

	c.PlayCard.PlayImage = playCardImage
}

func (c *Card) AddToHand(number int, widthFactor, heightFactor float64) {
	c.PlayCard.basePosX = float64(4 * number)
	c.PlayCard.basePosY = 11
	c.updatePlayCard(widthFactor, heightFactor)
}

func (pc *PlayCard) Click(x, y int) bool {
	if pc.CurrentPosX < float64(x) && pc.CurrentPosX+pc.CurrentWidth > float64(x) &&
		pc.CurrentPosY < float64(y) && pc.CurrentPosY+pc.currentHeight > float64(y) {
		return true
	}

	return false
}
