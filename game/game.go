package game

import (
	"bytes"
	"game/font"
	"game/game/cards"
	"game/game/screen"
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

	Screen screen.Screen
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
	menuHeader.Init(4, 2, 0, 0, true, "MENU", s, 2 /*text*/, white /*background*/, black /*border*/, white, screen.ScreenMain)
	menuCards := &entities.MenuItem{}
	menuCards.Init(4, 2, 0, 2, true, "CARDS", s, 1.5 /*text*/, white /*background*/, black /*border*/, blue, screen.ScreenCards)
	menuTech := &entities.MenuItem{}
	menuTech.Init(4, 2, 0, 4, true, "TECH", s, 1.5 /*text*/, white /*background*/, black /*border*/, green, screen.ScreenTech)
	menuAnna := &entities.MenuItem{}
	menuAnna.Init(4, 2, 0, 6, true, "ANNA", s, 1.5 /*text*/, pink /*background*/, black /*border*/, pink, screen.ScreenAnna)
	menuSettings := &entities.MenuItem{}
	menuSettings.Init(4, 2, 0, 14, true, "SETTINGS", s, 1 /*text*/, grey /*background*/, black /*border*/, grey, screen.ScreenSettings)
	g.MenuItems = append(g.MenuItems, menuHeader, menuCards, menuTech, menuAnna, menuSettings)

	menuDivider := &entities.Divider{}
	menuDivider.Init(true, 4, 16, 12, entities.Left, black)
	menuDividerMiddle := &entities.Divider{}
	menuDividerMiddle.Init(true, 4, 16, 9, entities.Left, darkGrey)
	menuDividerRight := &entities.Divider{}
	menuDividerRight.Init(true, 4, 16, 3, entities.Left, black)
	g.Dividers = append(g.Dividers, menuDivider, menuDividerMiddle, menuDividerRight)

	card1 := &cards.Card{}
	card1.Init(4, 8, 4, 0, "Test", "Does something\ncool", s, 2, 1, white, black, red)
	card2 := &cards.Card{}
	card2.Init(4, 8, 8, 0, "Test", "Does something\ncool", s, 2, 1, white, black, red)
	card3 := &cards.Card{}
	card3.Init(4, 8, 12, 0, "Test", "Does something\ncool", s, 2, 1, white, black, red)
	card4 := &cards.Card{}
	card4.Init(4, 8, 4, 8, "Test", "Does something\ncool", s, 2, 1, white, black, red)
	card5 := &cards.Card{}
	card5.Init(4, 8, 8, 8, "Test", "Does something\ncool", s, 2, 1, white, black, red)
	card6 := &cards.Card{}
	card6.Init(4, 8, 12, 8, "Test", "Does something\ncool", s, 2, 1, white, black, red)
	g.Cards = append(g.Cards, card1, card2, card3, card4, card5, card6)

	return g
}

func (g *Game) Update() error {
	//timeDelta := float64(time.Since(g.PreviousUpdateTime).Milliseconds())
	g.PreviousUpdateTime = time.Now()

	g.touchIDs = ebiten.AppendTouchIDs(g.touchIDs[:0])

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for _, menuItem := range g.MenuItems {
			newScreen := menuItem.Click(x, y)
			switch newScreen {
			case screen.ScreenInvalid:
				continue
			default:
				g.Screen = newScreen
			}
		}
	}

	return nil
}
