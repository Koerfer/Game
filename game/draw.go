package game

import (
	"fmt"
	"game/game/cards"
	screen2 "game/game/screen"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"math"
)

func (g *Game) Draw(screen *ebiten.Image) {
	g.CalculateNewFactors()

	switch g.Screen {
	case screen2.ScreenCards:
		g.DrawCards(screen)
	case screen2.ScreenMain:
		g.DrawMain(screen)
	case screen2.ScreenPlay:
		g.DrawPlay(screen)
	default:
		// nothing
	}

	g.DrawDividers(screen)
	g.DrawMenuItems(screen)
}

func (g *Game) CalculateNewFactors() {
	if g.WindowSize.Changed() {
		if g.WindowSize.CalculateNewFactorAndCheckIfChanged() {
			for _, menuItem := range g.MenuItems {
				menuItem.UpdateSize(g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)
			}
			for _, card := range g.Cards.Cards {
				card.Update(g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)
			}
			for _, divider := range g.MainDividers {
				divider.UpdateSize(g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)
			}
			for _, divider := range g.PlayDividers {
				divider.UpdateSize(g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)
			}
			g.StartButton.UpdateSize(g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)

			g.WindowSize.PreviousHeightFactor = g.WindowSize.CurrentHeightFactor
			g.WindowSize.PreviousWidthFactor = g.WindowSize.CurrentWidthFactor
		}

		g.WindowSize.PreviousScreenHeight = g.WindowSize.CurrentScreenHeight
		g.WindowSize.PreviousScreenWidth = g.WindowSize.CurrentScreenWidth
	}
}

func (g *Game) DrawDividers(screen *ebiten.Image) {
	for _, divider := range g.MainDividers {
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(float64(divider.CurrentPosX), float64(divider.CurrentPosY))
		screen.DrawImage(divider.Image, &ebiten.DrawImageOptions{
			GeoM: g.op.GeoM,
		})
	}
}

func (g *Game) DrawMenuItems(screen *ebiten.Image) {
	for _, menuItem := range g.MenuItems {
		if !menuItem.Shown {
			continue
		}

		g.op.GeoM.Reset()
		g.op.GeoM.Translate(menuItem.CurrentPosX, menuItem.CurrentPosY)
		screen.DrawImage(menuItem.Image, &ebiten.DrawImageOptions{
			GeoM: g.op.GeoM,
		})
	}
}

func (g *Game) DrawCards(screen *ebiten.Image) {
	for _, card := range g.Cards.Cards {
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(card.CurrentPosX, card.CurrentPosY)
		switch card.State {
		case cards.StateLocked:
			screen.DrawImage(card.LockedImage, &ebiten.DrawImageOptions{
				GeoM: g.op.GeoM,
			})
		case cards.StateUnlocked:
			screen.DrawImage(card.UnlockedImage, &ebiten.DrawImageOptions{
				GeoM: g.op.GeoM,
			})
		case cards.StateSelected:
			screen.DrawImage(card.SelectedImage, &ebiten.DrawImageOptions{
				GeoM: g.op.GeoM,
			})
		}
	}
}

func (g *Game) DrawMain(screen *ebiten.Image) {
	if g.StartButton.Shown {
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(g.StartButton.CurrentPosX, g.StartButton.CurrentPosY)
		screen.DrawImage(g.StartButton.Image, &ebiten.DrawImageOptions{
			GeoM: g.op.GeoM,
		})
	}
}

func (g *Game) DrawPlay(screen *ebiten.Image) {
	timeRemainingText := fmt.Sprintf(`Time: %02d:%02d`, int(math.Max(g.PlayState.TimeRemaining.Minutes(), 0)), int(math.Max(g.PlayState.TimeRemaining.Seconds(), 0))%60)
	waveText := fmt.Sprintf(`Wave: %d`, g.PlayState.Wave)
	numberMonstersText := fmt.Sprintf(`#Monsters: %d`, g.PlayState.NumberOfMonsters)
	hpPerMonsterText := fmt.Sprintf(`HP Per Monster: %d`, int(g.PlayState.HPPerMonster))
	if g.PlayState.HPPerMonster < 100 {
		hpPerMonsterText = fmt.Sprintf(`HP Per Monster: %.1f`, g.PlayState.HPPerMonster)
	}
	numberMonstersRemainingText := fmt.Sprintf(`Monsters Remaining: %d`, g.PlayState.MonstersRemaining)
	monstersAttackedText := fmt.Sprintf(`Monsters Attacked: %d`, g.PlayState.NumberOfMonstersAttacked)
	damagePerSecondSingleText := fmt.Sprintf(`Single Damage Per Second: %d`, int(g.PlayState.DamagePerSecond*g.PlayState.SingleTargetBoost))
	if g.PlayState.DamagePerSecond*g.PlayState.SingleTargetBoost < 100 {
		damagePerSecondSingleText = fmt.Sprintf(`Single Damage Per Second: %.1f`, g.PlayState.DamagePerSecond*g.PlayState.SingleTargetBoost)
	}
	damagePerSecondMultiText := fmt.Sprintf(`Multi Damage Per Second: %d`, int(g.PlayState.DamagePerSecond))
	if g.PlayState.DamagePerSecond < 100 {
		damagePerSecondMultiText = fmt.Sprintf(`Multi Damage Per Second: %.1f`, g.PlayState.DamagePerSecond)
	}
	cardUpgradesText := fmt.Sprintf(`Card Upgrades Available: %d`, g.Cards.Upgrades)

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(4.3*g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor*0.5)
	text.Draw(screen, waveText, &text.GoTextFace{
		Source: g.font,
		Size:   g.WindowSize.CurrentHeightFactor * 1,
	}, op)

	op.GeoM.Translate(4*g.WindowSize.CurrentWidthFactor, 0)
	text.Draw(screen, numberMonstersText, &text.GoTextFace{
		Source: g.font,
		Size:   g.WindowSize.CurrentHeightFactor * 1,
	}, op)

	op.GeoM.Translate(-4*g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor*1.75)
	text.Draw(screen, timeRemainingText, &text.GoTextFace{
		Source: g.font,
		Size:   g.WindowSize.CurrentHeightFactor * 1,
	}, op)

	op.GeoM.Translate(4*g.WindowSize.CurrentWidthFactor, 0)
	text.Draw(screen, hpPerMonsterText, &text.GoTextFace{
		Source: g.font,
		Size:   g.WindowSize.CurrentHeightFactor * 1,
	}, op)

	if g.PlayState.Playing {
		op.GeoM.Translate(-4*g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor*2.5)
		text.Draw(screen, numberMonstersRemainingText, &text.GoTextFace{
			Source: g.font,
			Size:   g.WindowSize.CurrentHeightFactor * 1,
		}, op)

		op.GeoM.Translate(0, g.WindowSize.CurrentHeightFactor*1.1)
		text.Draw(screen, monstersAttackedText, &text.GoTextFace{
			Source: g.font,
			Size:   g.WindowSize.CurrentHeightFactor * 1,
		}, op)

		op.GeoM.Translate(0, g.WindowSize.CurrentHeightFactor*1.1)
		text.Draw(screen, damagePerSecondMultiText, &text.GoTextFace{
			Source: g.font,
			Size:   g.WindowSize.CurrentHeightFactor * 1,
		}, op)

		op.GeoM.Translate(0, g.WindowSize.CurrentHeightFactor*1.1)
		text.Draw(screen, damagePerSecondSingleText, &text.GoTextFace{
			Source: g.font,
			Size:   g.WindowSize.CurrentHeightFactor * 1,
		}, op)
	} else {
		op.GeoM.Translate(-4*g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor*2.5)
		text.Draw(screen, cardUpgradesText, &text.GoTextFace{
			Source: g.font,
			Size:   g.WindowSize.CurrentHeightFactor * 1,
		}, op)
	}

	for _, card := range g.Cards.Selected {
		if card == nil {
			continue
		}
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(card.CurrentPosX, card.CurrentPosY)
		screen.DrawImage(card.PlayImage, &ebiten.DrawImageOptions{
			GeoM: g.op.GeoM,
		})
	}

	for _, divider := range g.PlayDividers {
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(float64(divider.CurrentPosX), float64(divider.CurrentPosY))
		screen.DrawImage(divider.Image, &ebiten.DrawImageOptions{
			GeoM: g.op.GeoM,
		})
	}

	for _, timer := range g.Timers {
		if timer == nil {
			continue
		}
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(timer.CurrentPosX, timer.CurrentPosY)
		screen.DrawImage(timer.Image, &ebiten.DrawImageOptions{
			GeoM: g.op.GeoM,
		})
	}
}
