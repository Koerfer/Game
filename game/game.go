package game

import (
	"fmt"
	"game/game/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"time"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

type Game struct {
	touchIDs []ebiten.TouchID
	op       ebiten.DrawImageOptions

	Squares            []*entities.Square
	Gravity            float64
	PreviousUpdateTime time.Time
}

func (g *Game) Update() error {
	timeDelta := float64(time.Since(g.PreviousUpdateTime).Milliseconds())
	g.PreviousUpdateTime = time.Now()

	g.touchIDs = ebiten.AppendTouchIDs(g.touchIDs[:0])

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for _, square := range g.Squares {
			square.Click(float64(x), float64(y))
		}
	}

	for _, square := range g.Squares {
		square.Update(ScreenWidth, ScreenHeight, timeDelta, g.Gravity)
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, square := range g.Squares {
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(square.PosX, square.PosY)

		screen.DrawImage(square.Image, &ebiten.DrawImageOptions{
			GeoM:          g.op.GeoM,
			CompositeMode: ebiten.CompositeModeSourceOver,
		})
	}
	msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f`, ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}
