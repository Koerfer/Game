package game

import (
	screen2 "game/game/screen"
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Draw(screen *ebiten.Image) {
	g.CalculateNewFactors()

	g.DrawDividers(screen)
	g.DrawMenuItems(screen)

	if g.Screen == screen2.ScreenCards {
		g.DrawCards(screen)
	}

	g.PrintDebug(screen)
}

func (g *Game) CalculateNewFactors() {
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
}

func (g *Game) DrawDividers(screen *ebiten.Image) {
	for _, divider := range g.Dividers {
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
	for _, card := range g.Cards {
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(card.CurrentPosX, card.CurrentPosY)
		screen.DrawImage(card.Image, &ebiten.DrawImageOptions{
			GeoM: g.op.GeoM,
		})
	}
}

func (g *Game) PrintDebug(screen *ebiten.Image) {
	for _, divider := range g.Dividers {
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(float64(divider.CurrentPosX), float64(divider.CurrentPosY))
		screen.DrawImage(divider.Image, &ebiten.DrawImageOptions{
			GeoM: g.op.GeoM,
		})
	}
}
