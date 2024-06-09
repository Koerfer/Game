package entities

import (
	"game/game/helper"
	"game/game/screen"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"image/color"
	"math"
)

type Button struct {
	baseWidth    float64
	baseHeight   float64
	basePosX     float64
	basePosY     float64
	BaseTextSize float64

	CurrentWidth    float64
	CurrentHeight   float64
	CurrentPosX     float64
	CurrentPosY     float64
	CurrentTextSize float64

	Shown bool

	Image            *ebiten.Image
	BackgroundColour color.Color
	Colour           color.Color
	TextColour       color.Color
	Name             string
	Font             *text.GoTextFaceSource

	Screen screen.Screen
}

func (b *Button) Init(width, height, x, y float64, shown bool, name string, font *text.GoTextFaceSource, textSize float64, textColour, backgroundColour, colour color.Color, screen screen.Screen) {
	buttonImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: int(x), Y: int(y)},
		Max: image.Point{X: int(x + width), Y: int(y + height)},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	b.baseWidth = width
	b.baseHeight = height
	b.basePosX = x
	b.basePosY = y

	b.CurrentWidth = width
	b.CurrentHeight = height
	b.CurrentPosX = x
	b.CurrentPosY = y

	b.Shown = shown

	b.Image = buttonImage
	b.BackgroundColour = backgroundColour
	b.Colour = colour
	b.TextColour = textColour
	b.Name = name
	b.Font = font
	b.BaseTextSize = textSize

	b.Screen = screen
}

func (b *Button) Update() {
	if b.CurrentTextSize*float64(len(b.Name)) > b.CurrentWidth {
		b.CurrentTextSize = math.Min(b.CurrentTextSize, 1.4*b.CurrentWidth/float64(len(b.Name)))
	}

	buttonImage := ebiten.NewImageWithOptions(image.Rectangle{
		Min: image.Point{X: int(b.CurrentPosX), Y: int(b.CurrentPosY)},
		Max: image.Point{X: int(b.CurrentPosX + b.CurrentWidth), Y: int(b.CurrentPosY + b.CurrentHeight)},
	}, &ebiten.NewImageOptions{Unmanaged: false})

	buttonImage.Fill(b.BackgroundColour)
	for i := 0; i <= int(b.CurrentWidth+b.CurrentPosX); i++ {
		for n := 0; n <= 5; n++ {
			buttonImage.Set(i, int(b.CurrentPosY)+n, b.Colour)
			buttonImage.Set(i, int(b.CurrentPosY+b.CurrentHeight)-n, b.Colour)
		}
	}
	for i := 0; i <= int(b.CurrentHeight+b.CurrentPosY); i++ {
		for n := 0; n <= 5; n++ {
			buttonImage.Set(int(b.CurrentPosX)+n, i, b.Colour)
			buttonImage.Set(int(b.CurrentPosX+b.CurrentWidth)-n, i, b.Colour)
		}
	}

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(b.TextColour)
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter

	middleX := b.CurrentWidth/2 + b.CurrentPosX
	middleY := b.CurrentHeight/2 + b.CurrentPosY
	op.GeoM.Translate(middleX, middleY)

	text.Draw(buttonImage, b.Name, &text.GoTextFace{
		Source: b.Font,
		Size:   b.CurrentTextSize,
	}, op)

	b.Image = buttonImage
}

func (b *Button) UpdateSize(widthFactor, heightFactor float64) {
	newWidth, newHeight, newX, newY := helper.GetNewSizeAndPosition(b.baseWidth, b.baseHeight, b.basePosX, b.basePosY, widthFactor, heightFactor, 0, 0)
	newTextSize := helper.GetNewTextSize(b.BaseTextSize, heightFactor, newWidth, b.Name)

	b.CurrentWidth = newWidth
	b.CurrentHeight = newHeight
	b.CurrentPosX = newX
	b.CurrentPosY = newY
	b.CurrentTextSize = newTextSize

	b.Update()
}

func (b *Button) Click(x, y int) screen.Screen {
	if !b.Shown {
		return screen.ScreenNothing
	}

	if b.CurrentPosX < float64(x) && b.CurrentPosX+b.CurrentWidth > float64(x) &&
		b.CurrentPosY < float64(y) && b.CurrentPosY+b.CurrentHeight > float64(y) {
		return b.Screen
	}

	return screen.ScreenNothing
}
