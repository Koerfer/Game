package cards

import (
	"fmt"
	"game/game/colours"
	"game/game/helper"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"time"
)

type PlayCard struct {
	Id int

	baseWidth  float64
	baseHeight float64
	basePosX   float64
	basePosY   float64

	currentWidth  float64
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
	passiveMultiBoost, activeMultiBoost int, activationTime, coolDown time.Duration) {
	pc.BaseDescription = description

	if passiveDamageBoost != 1 {
		pc.Description = fmt.Sprintf(description, int(passiveDamageBoost), int(activeSingleTargetDamageBoost), int(coolDown.Seconds()))
	} else if passiveMultiBoost != 1 {
		pc.Description = fmt.Sprintf(description, passiveMultiBoost, activeMultiBoost, int(coolDown.Seconds()))
	} else {
		pc.Description = description
	}

	pc.PassiveDamageBoost = passiveDamageBoost
	pc.ActiveSingleTargetDamageBoost = activeSingleTargetDamageBoost

	pc.PassiveMultiTargetBoost = passiveMultiBoost
	pc.ActiveMultiTargetBoost = activeMultiBoost

	pc.ActiveTime = activationTime
	pc.CoolDown = coolDown
	pc.CoolDownRemaining = 0
}

func (c *Card) updatePlayCard(widthFactor, heightFactor float64) {
	if c.PlayCard.PassiveDamageBoost != 1 {
		c.PlayCard.Description = fmt.Sprintf(c.PlayCard.BaseDescription, int(c.PlayCard.PassiveDamageBoost), int(c.PlayCard.ActiveSingleTargetDamageBoost), int(c.PlayCard.CoolDown.Seconds()))
	} else if c.PlayCard.PassiveMultiTargetBoost != 1 {
		c.PlayCard.Description = fmt.Sprintf(c.PlayCard.BaseDescription, c.PlayCard.PassiveMultiTargetBoost, c.PlayCard.ActiveMultiTargetBoost, int(c.PlayCard.CoolDown.Seconds()))
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

	c.PlayCard.currentWidth = newWidth
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
	if pc.CurrentPosX < float64(x) && pc.CurrentPosX+pc.currentWidth > float64(x) &&
		pc.CurrentPosY < float64(y) && pc.CurrentPosY+pc.currentHeight > float64(y) {
		return true
	}

	return false
}

func (pc *PlayCard) TimerImage(currentTime, totalTime time.Duration, active bool, verticalScale float64) *Timer {
	var colour color.Color
	if active {
		colour = colours.Green
	} else {
		colour = colours.Red
	}

	xMin := pc.CurrentPosX
	xMax := xMin + pc.currentWidth*float64(currentTime.Milliseconds())/float64(totalTime.Milliseconds())
	yMin := pc.CurrentPosY - 10*verticalScale/50
	yMax := pc.CurrentPosY

	if int(pc.currentWidth*float64(currentTime.Milliseconds())/float64(totalTime.Milliseconds())) <= 0 {
		return nil
	}

	timerBar := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: int(xMin), Y: int(yMin)},
		Max: image.Point{X: int(xMax), Y: int(yMax)},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	timerBar.Fill(colour)

	return &Timer{
		CurrentPosX: pc.CurrentPosX,
		CurrentPosY: pc.CurrentPosY - 10*verticalScale/50,
		Image:       timerBar,
	}
}

type Timer struct {
	CurrentPosX float64
	CurrentPosY float64

	Image *ebiten.Image
}
