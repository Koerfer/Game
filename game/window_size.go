package game

type WindowSize struct {
	PreviousScreenWidth  int
	PreviousScreenHeight int

	CurrentScreenWidth  int
	CurrentScreenHeight int

	PreviousWidthFactor  int
	PreviousHeightFactor int

	CurrentWidthFactor  int
	CurrentHeightFactor int
}

func (ws *WindowSize) Changed() bool {
	if ws.PreviousScreenHeight != ws.CurrentScreenHeight {
		return true
	}
	if ws.PreviousScreenWidth != ws.CurrentScreenWidth {
		return true
	}

	return false
}

func (ws *WindowSize) CalculateNewFactorAndCheckIfChanged() bool {
	ws.CurrentWidthFactor = ws.CurrentScreenWidth / 18
	ws.CurrentHeightFactor = ws.CurrentScreenHeight / 18

	if ws.PreviousWidthFactor != ws.CurrentWidthFactor {
		return true
	}
	if ws.PreviousHeightFactor != ws.CurrentHeightFactor {
		return true
	}

	return false
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.WindowSize.CurrentScreenWidth = outsideWidth
	g.WindowSize.CurrentScreenHeight = outsideHeight

	return outsideWidth, outsideHeight
}
