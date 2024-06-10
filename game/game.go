package game

import (
	"game/game/cards"
	"game/game/colours"
	"game/game/font"
	"game/game/play"
	"game/game/screen"
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
	PreviousSaveTime   time.Time
	MenuItems          []*entities.Button
	MainDividers       []*entities.Divider
	PlayDividers       []*entities.Divider
	Cards              *cards.Cards
	StartButton        *entities.Button
	Timers             []*cards.Timer

	PlayState   *play.State
	HighestWave int

	Screen screen.Screen
}

func Init() *Game {
	g := &Game{}
	ebiten.SetTPS(60)
	g.PreviousUpdateTime = time.Now()
	g.PreviousSaveTime = time.Now()

	saveState := g.Load()
	g.newInit()

	if saveState == nil {
		return g
	}

	g.UpdateToMatchLoadedState(saveState)

	return g
}

func (g *Game) newInit() {
	screenWidth, screenHeight := ebiten.WindowSize()
	g.WindowSize = &WindowSize{
		CurrentScreenWidth:  screenWidth,
		CurrentScreenHeight: screenHeight,
	}

	g.Timers = make([]*cards.Timer, 3)

	g.font = font.GetBold()

	menuHeader := &entities.Button{}
	menuHeader.Init(4, 2, 0, 0, true, "MENU", g.font, 2 /*text*/, colours.White /*background*/, colours.Black /*border*/, colours.White, screen.ScreenNothing)
	menuPlay := &entities.Button{}
	menuPlay.Init(4, 2, 0, 2, true, "PLAY", g.font, 1.5 /*text*/, colours.White /*background*/, colours.Black /*border*/, colours.Green, screen.ScreenMain)
	menuCards := &entities.Button{}
	menuCards.Init(4, 2, 0, 4, false, "CARDS", g.font, 1.5 /*text*/, colours.White /*background*/, colours.Black /*border*/, colours.Blue, screen.ScreenCards)
	menuTech := &entities.Button{}
	menuTech.Init(4, 2, 0, 6, false, "TECH", g.font, 1.5 /*text*/, colours.White /*background*/, colours.Black /*border*/, colours.Red, screen.ScreenTech)
	menuAnna := &entities.Button{}
	menuAnna.Init(4, 2, 0, 8, false, "ANNA", g.font, 1.5 /*text*/, colours.Pink /*background*/, colours.Black /*border*/, colours.Pink, screen.ScreenAnna)

	menuSave := &entities.Button{}
	menuSave.Init(4, 2, 0, 12, true, "SAVE", g.font, 1 /*text*/, colours.Green /*background*/, colours.DarkGreen /*border*/, colours.Green, screen.ScreenSave)
	menuSettings := &entities.Button{}
	menuSettings.Init(4, 2, 0, 14, true, "SETTINGS", g.font, 1 /*text*/, colours.Grey /*background*/, colours.DarkGrey /*border*/, colours.Grey, screen.ScreenSettings)
	g.MenuItems = append(g.MenuItems, menuHeader, menuPlay, menuCards, menuTech, menuAnna, menuSave, menuSettings)

	g.StartButton = &entities.Button{}
	g.StartButton.Init(8, 4, 6, 6, true, "START", g.font, 3 /*text*/, colours.Green /*background*/, colours.DarkGreen /*border*/, colours.Green, screen.ScreenPlay)

	menuDivider := &entities.Divider{}
	menuDivider.Init(true, 4, 0, 16, 12, entities.Left, colours.Black)
	menuDividerMiddle := &entities.Divider{}
	menuDividerMiddle.Init(true, 4, 0, 16, 9, entities.Left, colours.DarkGrey)
	menuDividerRight := &entities.Divider{}
	menuDividerRight.Init(true, 4, 0, 16, 3, entities.Left, colours.Black)
	g.MainDividers = append(g.MainDividers, menuDivider, menuDividerMiddle, menuDividerRight)

	waveDivider := &entities.Divider{}
	waveDivider.Init(true, 8, 0, 4, 5, entities.Left, colours.White)
	topDivider := &entities.Divider{}
	topDivider.Init(false, 4, 4, 12, 5, entities.Left, colours.White)
	g.PlayDividers = append(g.PlayDividers, waveDivider, topDivider)

	cs := &cards.Cards{}
	cs.Init()
	g.Cards = cs

	g.Screen = screen.ScreenNothing
}

func (g *Game) Update() error {
	timeDelta := time.Since(g.PreviousUpdateTime) * 10
	g.PreviousUpdateTime = time.Now()

	if time.Since(g.PreviousSaveTime) > 2*time.Minute {
		g.Save()
	}

	if g.Screen == screen.ScreenPlay {
		for i, selectedCard := range g.Cards.Selected {
			if selectedCard == nil {
				continue
			}
			if selectedCard.Active {
				selectedCard.ActiveRemaining -= timeDelta
				g.Timers[i] = selectedCard.TimerImage(selectedCard.ActiveRemaining, selectedCard.ActiveTime, true, g.WindowSize.CurrentHeightFactor)
				continue
			}
			if selectedCard.CoolDownRemaining > 0 {
				selectedCard.CoolDownRemaining -= timeDelta
				g.Timers[i] = selectedCard.TimerImage(selectedCard.CoolDownRemaining, selectedCard.CoolDown, false, g.WindowSize.CurrentHeightFactor)
				if selectedCard.CoolDownRemaining <= 0 {
					selectedCard.CoolDownRemaining = 0
					g.Timers[i] = nil
				}
			}
		}

		pause := g.PlayState.Update(timeDelta, g.HighestWave)
		switch pause {
		case -1:
			for _, selectedCard := range g.Cards.Selected {
				if selectedCard == nil {
					continue
				}
				selectedCard.CoolDownRemaining = 0
				selectedCard.ActiveRemaining = 0
				selectedCard.Active = false
			}
			g.StartButton.Name = "START"
			g.Timers = make([]*cards.Timer, 3)
			g.StartButton.Update()
			g.PlayState = nil
			g.Screen = screen.ScreenMain
		case 0:
			// do nothing
		case 1:
			for _, card := range g.Cards.Cards {
				if card.State == cards.StateLocked {
					card.State = cards.StateUnlocked
					break
				}
			}
		case 2:
			g.Cards.Upgrades += 1
		default:
			log.Fatal("invalid pause reason")
		}

		if g.PlayState != nil && !g.PlayState.Playing && g.HighestWave < g.PlayState.Wave-1 {
			g.HighestWave = g.PlayState.Wave - 1
			if g.HighestWave > 5 {
				g.MenuItems[2].Shown = true
			}
		}
	}

	g.touchIDs = ebiten.AppendTouchIDs(g.touchIDs[:0])

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for _, menuItem := range g.MenuItems {
			newScreen := menuItem.Click(x, y)
			switch newScreen {
			case screen.ScreenNothing:
				continue
			case screen.ScreenSave:
				g.Save()
				g.PreviousSaveTime = time.Now()
				continue
			default:
				g.Screen = newScreen
			}
		}

		switch g.Screen {
		case screen.ScreenCards:
			if g.PlayState != nil && g.PlayState.Playing {
				break
			}
			for _, card := range g.Cards.Cards {
				if ebiten.IsKeyPressed(ebiten.KeyShiftRight) || ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShift) {
					upgraded := card.Upgrade(x, y, g.Cards.Upgrades)
					if upgraded {
						g.Cards.Upgrades -= 1
						card.Update(g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)
					}
					continue
				}

				added := card.Click(x, y, g.Cards.NumberSelected)

				switch added {
				case 1:
					for i, selectedCard := range g.Cards.Selected {
						if selectedCard == nil {
							card.AddToHand(i+1, g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)
							g.Cards.Selected[i] = card.PlayCard
							break
						}
					}
				case -1:
					for i, selectedCard := range g.Cards.Selected {
						if selectedCard == nil {
							continue
						}
						if card.Id == selectedCard.Id {
							g.Cards.Selected[i] = nil
							break
						}
					}
				default:
					// do nothing
				}
				g.Cards.NumberSelected += added
			}
		case screen.ScreenMain:
			newScreen := g.StartButton.Click(x, y)
			switch newScreen {
			case screen.ScreenPlay:
				g.Screen = newScreen
				if g.PlayState == nil {
					g.PlayState = play.Start(g.Cards.Selected)
					g.StartButton.Name = "CONTINUE"
					g.StartButton.Update()
				}
			default:
				// do nothing
			}

		case screen.ScreenPlay:
			for i, selectedCard := range g.Cards.Selected {
				if !g.PlayState.Playing {
					continue
				}
				if selectedCard == nil {
					continue
				}
				if selectedCard.Active {
					continue
				}
				if selectedCard.CoolDownRemaining > 0 {
					continue
				}
				if selectedCard.Click(x, y) {
					g.PlayState.CardActivation(selectedCard, i)
				}
			}
		default:
			// do nothing
		}

	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Screen = screen.ScreenMain
	}

	return nil
}
