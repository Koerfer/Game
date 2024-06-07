package game

import (
	"bytes"
	"fmt"
	"game/font"
	"game/game/cards"
	"image/color"
	"log"
	"time"

	"game/game/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Game struct {
	touchIDs   []ebiten.TouchID
	op         ebiten.DrawImageOptions
	font       *text.GoTextFaceSource
	WindowSize *WindowSize

	PreviousUpdateTime time.Time
	MenuItems          []*entities.MenuItem
	Dividers           []*entities.Divider
	Cards              []*cards.Card
}

func Init() *Game {
	g := &Game{}
	ebiten.SetTPS(120)
	g.PreviousUpdateTime = time.Now()
	s, err := text.NewGoTextFaceSource(bytes.NewReader(font.MonoBold))
	if err != nil {
		log.Fatal(err)
	}
	g.font = s
	screenWidth, screenHeight := ebiten.WindowSize()
	g.WindowSize = &WindowSize{
		CurrentScreenWidth:  screenWidth,
		CurrentScreenHeight: screenHeight,
	}

	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	grey := color.RGBA{R: 150, G: 150, B: 150, A: 255}
	darkGrey := color.RGBA{R: 70, G: 70, B: 70, A: 255}
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}
	pink := color.RGBA{R: 255, G: 150, B: 200, A: 255}
	menuHeader := &entities.MenuItem{}
	menuHeader.Init(4, 2, 0, 0, true, "MENU", s, 2 /*text*/, white /*background*/, black /*border*/, white, nil)
	menuCards := &entities.MenuItem{}
	menuCards.Init(4, 2, 0, 2, true, "CARDS", s, 1.5 /*text*/, white /*background*/, black /*border*/, blue, nil)
	menuTech := &entities.MenuItem{}
	menuTech.Init(4, 2, 0, 4, true, "TECH", s, 1.5 /*text*/, white /*background*/, black /*border*/, green, nil)
	menuAnna := &entities.MenuItem{}
	menuAnna.Init(4, 2, 0, 6, false, "ANNA", s, 1.5 /*text*/, pink /*background*/, black /*border*/, pink, nil)
	menuSettings := &entities.MenuItem{}
	menuSettings.Init(4, 2, 0, 14, true, "SETTINGS", s, 1 /*text*/, grey /*background*/, black /*border*/, grey, nil)
	g.MenuItems = append(g.MenuItems, menuHeader, menuCards, menuTech, menuAnna, menuSettings)

	menuDivider := &entities.Divider{}
	menuDivider.Init(true, 4, 16, 12, entities.Left, black)
	menuDividerMiddle := &entities.Divider{}
	menuDividerMiddle.Init(true, 4, 16, 9, entities.Left, darkGrey)
	menuDividerRight := &entities.Divider{}
	menuDividerRight.Init(true, 4, 16, 3, entities.Left, black)
	g.Dividers = append(g.Dividers, menuDivider, menuDividerMiddle, menuDividerRight)

	card := &cards.Card{}
	card.Init(4, 8, 4, 0, "Test", "Does something \n cool", s, 2, 1, white, black, red)
	card2 := &cards.Card{}
	card2.Init(4, 8, 8, 0, "Test", "Does something \n cool", s, 2, 1, white, black, red)
	card3 := &cards.Card{}
	card3.Init(4, 8, 12, 0, "Test", "Does something \n cool", s, 2, 1, white, black, red)
	card4 := &cards.Card{}
	card4.Init(4, 8, 4, 8, "Test", "Does something \n cool", s, 2, 1, white, black, red)
	g.Cards = append(g.Cards, card, card2, card3, card4)

	return g
}

func (g *Game) Update() error {
	timeDelta := float64(time.Since(g.PreviousUpdateTime).Milliseconds())
	g.PreviousUpdateTime = time.Now()

	g.touchIDs = ebiten.AppendTouchIDs(g.touchIDs[:0])

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		_ = x + y + int(timeDelta)

	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.WindowSize.Changed() {
		if g.WindowSize.CalculateNewFactorAndCheckIfChanged() {
			for _, menuItem := range g.MenuItems {
				menuItem.UpdateSize(g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)
			}
			for _, card := range g.Cards {
				card.UpdateSize(g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)
			}
			for _, divider := range g.Dividers {
				divider.UpdateSize(g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)
			}

			g.WindowSize.PreviousHeightFactor = g.WindowSize.CurrentHeightFactor
			g.WindowSize.PreviousWidthFactor = g.WindowSize.CurrentWidthFactor
		}

		g.WindowSize.PreviousScreenHeight = g.WindowSize.CurrentScreenHeight
		g.WindowSize.PreviousScreenWidth = g.WindowSize.CurrentScreenWidth
	}

	for _, divider := range g.Dividers {
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(float64(divider.CurrentPosX), float64(divider.CurrentPosY))
		screen.DrawImage(divider.Image, &ebiten.DrawImageOptions{
			GeoM: g.op.GeoM,
		})
	}

	for _, menuItem := range g.MenuItems {
		if !menuItem.Shown {
			continue
		}

		g.op.GeoM.Reset()
		g.op.GeoM.Translate(float64(menuItem.CurrentPosX), float64(menuItem.CurrentPosY))
		screen.DrawImage(menuItem.Image, &ebiten.DrawImageOptions{
			GeoM: g.op.GeoM,
		})
	}

	for _, card := range g.Cards {
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(float64(card.CurrentPosX), float64(card.CurrentPosY))
		screen.DrawImage(card.Image, &ebiten.DrawImageOptions{
			GeoM: g.op.GeoM,
		})
	}

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)

	msg := fmt.Sprintf(`TPS: %0.2f - FPS: %0.2f`, ebiten.ActualTPS(), ebiten.ActualFPS())
	op.GeoM.Translate(float64(g.WindowSize.PreviousScreenWidth)-170, 0)
	text.Draw(screen, msg, &text.GoTextFace{
		Source: g.font,
		Size:   12,
	}, op)
	op.GeoM.Reset()
}
