package game

import (
	"bytes"
	"fmt"
	"game/font"
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
		PreviousScreenWidth:  18,
		PreviousScreenHeight: 18,
		CurrentScreenWidth:   screenWidth,
		CurrentScreenHeight:  screenHeight,
	}

	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}
	menuHeader := &entities.MenuItem{}
	menuHeader.Init(4, 2, 0, 0, true, "MENU", s, 2 /*text*/, white /*background*/, black /*border*/, black, nil)
	menuCards := &entities.MenuItem{}
	menuCards.Init(4, 2, 0, 2, true, "CARDS", s, 1 /*text*/, white /*background*/, black /*border*/, blue, nil)
	menuTech := &entities.MenuItem{}
	menuTech.Init(4, 2, 0, 4, true, "TECH", s, 1 /*text*/, white /*background*/, black /*border*/, green, nil)
	g.MenuItems = append(g.MenuItems, menuHeader, menuCards, menuTech)

	menuDivider := &entities.Divider{}
	menuDivider.Init(true, 4, 50, color.RGBA{R: 255, G: 255, B: 255, A: 0})
	g.Dividers = append(g.Dividers, menuDivider)

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
