package game

import (
	"fmt"
	"game/game/cards"
	"game/game/colours"
	"game/game/helper"
	screen2 "game/game/screen"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"math"
	"strings"
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

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)
	msg := fmt.Sprintf(`TPS: %0.2f - FPS: %0.2f`, ebiten.ActualTPS(), ebiten.ActualFPS())
	op.GeoM.Translate(float64(g.WindowSize.CurrentScreenWidth)-170, 0)
	text.Draw(screen, msg, &text.GoTextFace{
		Source: g.font,
		Size:   12,
	}, op)
	op.GeoM.Reset()
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
			g.showUpgradeText(screen, card.CurrentPosX, card.CurrentPosY, card.CurrentWidth, card.CurrentHeight,
				helper.GetNewTextSize(card.BaseNameTextSize, g.WindowSize.CurrentHeightFactor, card.CurrentWidth, "UPGRADE"))
		case cards.StateSelected:
			screen.DrawImage(card.SelectedImage, &ebiten.DrawImageOptions{
				GeoM: g.op.GeoM,
			})
			g.showUpgradeText(screen, card.CurrentPosX, card.CurrentPosY, card.CurrentWidth, card.CurrentHeight,
				helper.GetNewTextSize(card.BaseNameTextSize, g.WindowSize.CurrentHeightFactor, card.CurrentWidth, "UPGRADE"))
		}
	}
}

func (g *Game) showUpgradeText(screen *ebiten.Image, x, y, width, height, textSize float64) {
	if g.Cards.Upgrades == 0 {
		return
	}
	g.op.GeoM.Translate(6.5, height-textSize-10)
	background := ebiten.NewImage(int(width)-12, int(height-(height-textSize-10)-5))
	background.Fill(colours.GreenTrans)
	screen.DrawImage(background, &ebiten.DrawImageOptions{
		GeoM:  g.op.GeoM,
		Blend: ebiten.BlendSourceAtop,
	})

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(colours.Green)
	op.PrimaryAlign = text.AlignCenter

	posX := width/2 + x
	posY := y + height - textSize - 10
	op.GeoM.Translate(posX, posY)

	text.Draw(screen, "UPGRADE", &text.GoTextFace{
		Source: g.font,
		Size:   textSize,
	}, op)
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
	hpPerMonsterText := strings.Replace(fmt.Sprintf(`HP Per Monster: %.3g`, g.PlayState.HPPerMonster), "e+0", "e", 1)
	hpPerMonsterText = strings.Replace(hpPerMonsterText, "e+", "e", 1)
	numberMonstersRemainingText := fmt.Sprintf(`Monsters Remaining: %d`, g.PlayState.MonstersRemaining)
	monstersAttackedText := fmt.Sprintf(`Monsters Attacked: %d`, g.PlayState.NumberOfMonstersAttacked)
	damagePerSecondSingleText := strings.Replace(fmt.Sprintf(`Single Damage Per Second: %.3g`, g.PlayState.DamagePerSecond*g.PlayState.SingleTargetBoost), "e+0", "e", 1)
	damagePerSecondSingleText = strings.Replace(damagePerSecondSingleText, "e+", "e", 1)
	damagePerSecondMultiText := strings.Replace(fmt.Sprintf(`Multi Damage Per Second: %.3g`, g.PlayState.DamagePerSecond), "e+", "e", 1)
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

		op.GeoM.Translate(0, g.WindowSize.CurrentHeightFactor*1.5)
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

		if card.Active {
			width := int(card.CurrentWidth * float64(card.ActiveRemaining.Milliseconds()) / float64(card.ActiveTime.Milliseconds()))
			if width <= 0 {
				continue
			}
			g.op.GeoM.Translate(0, -g.WindowSize.CurrentHeightFactor/5)
			timerBar := ebiten.NewImage(width, int(g.WindowSize.CurrentHeightFactor/5))
			timerBar.Fill(colours.Green)
			screen.DrawImage(timerBar, &ebiten.DrawImageOptions{
				GeoM: g.op.GeoM,
			})
		} else if card.CoolDownRemaining != 0 {
			width := int(card.CurrentWidth * float64(card.CoolDownRemaining.Milliseconds()) / float64(card.CoolDown.Milliseconds()))
			if width <= 0 {
				continue
			}
			g.op.GeoM.Translate(0, -g.WindowSize.CurrentHeightFactor/5)
			timerBar := ebiten.NewImage(width, int(g.WindowSize.CurrentHeightFactor/5))
			timerBar.Fill(colours.Red)
			screen.DrawImage(timerBar, &ebiten.DrawImageOptions{
				GeoM: g.op.GeoM,
			})
		}

	}

	for _, divider := range g.PlayDividers {
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(float64(divider.CurrentPosX), float64(divider.CurrentPosY))
		screen.DrawImage(divider.Image, &ebiten.DrawImageOptions{
			GeoM: g.op.GeoM,
		})
	}
}
