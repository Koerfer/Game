package cards

import (
	"game/game/colours"
	"game/game/helper"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"math"
	"strings"
)

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func (c *Card) Update(widthFactor, heightFactor float64) {
	shiftX := 12.0
	if c.basePosX == 8 {
		shiftX = 8.0
	} else if c.basePosX == 12 {
		shiftX = 4.0
	}

	shiftY := 3.0
	if c.basePosY == 8 {
		shiftY = -1.0
	}

	newWidth, newHeight, newX, newY := helper.GetNewSizeAndPosition(c.baseWidth, c.baseHeight, c.basePosX, c.basePosY, widthFactor, heightFactor, shiftX, shiftY)
	newWidth -= 7
	newHeight -= 7
	newNameTextSize := helper.GetNewTextSize(c.BaseNameTextSize, heightFactor, newWidth, c.Name)

	c.updatePlayCard(widthFactor, heightFactor)
	c.CurrentWidth = newWidth
	c.CurrentHeight = newHeight
	c.CurrentPosX = newX
	c.CurrentPosY = newY
	c.LockedImage = c.createImage(newX, newY, newWidth, newHeight, newNameTextSize, heightFactor, StateLocked)
	c.UnlockedImage = c.createImage(newX, newY, newWidth, newHeight, newNameTextSize, heightFactor, StateUnlocked)
	c.SelectedImage = c.createImage(newX, newY, newWidth, newHeight, newNameTextSize, heightFactor, StateSelected)
}

func (c *Card) createImage(newX, newY, newWidth, newHeight, newNameTextSize, heightFactor float64, state State) *ebiten.Image {
	backgroundColour := c.BackgroundColour
	colour := c.Colour
	textColour := c.TextColour
	switch state {
	case StateLocked:
		backgroundColour = colours.DarkRed
		colour = colours.Red
	case StateSelected:
		backgroundColour = colours.DarkBlue
	case StateUnlocked:
		// do nothing
	}

	cardImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: int(newX), Y: int(newY)},
		Max: image.Point{X: int(newX + newWidth), Y: int(newY + newHeight)},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	cardImage.Fill(backgroundColour)
	for i := 0; i <= int(newWidth+newX); i++ {
		for n := 0; n <= 5; n++ {
			cardImage.Set(i, int(newY)+n, colour)
			cardImage.Set(i, int(newY+newHeight)-n, colour)
		}
	}
	for i := 0; i <= int(newHeight+newY); i++ {
		for n := 0; n <= 5; n++ {
			cardImage.Set(int(newX)+n, i, colour)
			cardImage.Set(int(newX+newWidth)-n, i, colour)
		}
	}

	switch state {
	case StateLocked:
		c.drawPadLock(cardImage, newWidth, newHeight, float32(newX), float32(newY))
	default:
		c.printName(cardImage, newNameTextSize, newWidth, newHeight, newX, newY, textColour)
		c.printDescription(cardImage, heightFactor, newWidth, newHeight, newX, newY, textColour)
	}

	return cardImage
}

func (c *Card) printName(cardImage *ebiten.Image, newNameTextSize, newWidth, newHeight, newX, newY float64, textColour color.Color) {
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(textColour)
	op.PrimaryAlign = text.AlignCenter

	posX := newWidth/2 + newX
	posY := newY + (newHeight*0.2-newNameTextSize)/2
	op.GeoM.Translate(posX, posY)

	text.Draw(cardImage, c.Name, &text.GoTextFace{
		Source: c.Font,
		Size:   newNameTextSize,
	}, op)
}

func (c *Card) printNamePlay(cardImage *ebiten.Image, newNameTextSize, newWidth, newHeight, newX, newY float64, textColour color.Color) {
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(textColour)
	op.PrimaryAlign = text.AlignCenter

	posX := newWidth/2 + newX
	posY := newY + (newHeight*0.3-newNameTextSize)/2
	op.GeoM.Translate(posX, posY)

	text.Draw(cardImage, c.Name, &text.GoTextFace{
		Source: c.Font,
		Size:   newNameTextSize,
	}, op)
}

func (c *Card) printDescription(cardImage *ebiten.Image, heightFactor, newWidth, newHeight, newX, newY float64, textColour color.Color) {
	newDescriptionTextSize := helper.GetNewTextSize(c.BaseDescriptionTextSize, heightFactor, newWidth, c.PlayCard.Description)
	widthOfText, _ := text.Measure(c.PlayCard.Description, &text.GoTextFace{
		Source: c.Font,
		Size:   newDescriptionTextSize,
	}, 0)

	splitDescription := strings.Split(c.PlayCard.Description, "\n")

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(textColour)
	op.PrimaryAlign = text.AlignCenter

	finalSize := newWidth / widthOfText * newDescriptionTextSize * 0.9
	posX := newWidth/2 + newX
	posY := newY + newHeight/2 - (finalSize*1.1*float64(len(splitDescription)))/6
	op.GeoM.Translate(posX, posY)
	for _, subDescription := range splitDescription {
		text.Draw(cardImage, subDescription, &text.GoTextFace{
			Source: c.Font,
			Size:   finalSize,
		}, op)
		op.GeoM.Translate(0, finalSize*1.1)
	}
}

func (c *Card) drawPadLock(cardImage *ebiten.Image, cardWidth, cardHeight float64, newX, newY float32) {
	size := float32(math.Min(cardWidth/2, cardHeight/1.7))

	var path vector.Path

	path.MoveTo(0, 0)
	path.LineTo(size, 0)
	path.LineTo(size, size*0.7)
	path.LineTo(0, size*0.7)
	path.LineTo(0, 0)
	path.LineTo(size/2*0.2, 0)
	path.Arc(size/2, -0.4*size/2, size/2*0.8, math.Pi, 0, vector.Clockwise)
	path.LineTo(size/2*1.8, 0)
	path.Close()

	var vs []ebiten.Vertex
	var is []uint16
	opStroke := &vector.StrokeOptions{}
	opStroke.Width = 30
	opStroke.LineJoin = vector.LineJoinRound
	vs, is = path.AppendVerticesAndIndicesForStroke(nil, nil, opStroke)

	for i := range vs {
		vs[i].DstX = vs[i].DstX + newX + (float32(cardWidth)-size)/2
		vs[i].DstY = vs[i].DstY + newY + (float32(cardHeight)-size*0.7)/2 + (0.4*size/2+size/2*0.8)/2
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = 1
		vs[i].ColorG = 0
		vs[i].ColorB = 0
		vs[i].ColorA = 1
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = true
	op.FillRule = ebiten.FillAll
	cardImage.DrawTriangles(vs, is, whiteSubImage, op)
}
