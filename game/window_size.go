package game

type WindowSize struct {
	PreviousScreenWidth  int
	PreviousScreenHeight int

	CurrentScreenWidth  int
	CurrentScreenHeight int

	PreviousWidthFactor  float64
	PreviousHeightFactor float64

	CurrentWidthFactor  float64
	CurrentHeightFactor float64
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
	ws.CurrentWidthFactor = float64(ws.CurrentScreenWidth) / 16
	ws.CurrentHeightFactor = float64(ws.CurrentScreenHeight) / 16

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
