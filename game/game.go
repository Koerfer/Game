package game

import (
	"fmt"
	"game/game/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"time"
)

type Game struct {
	touchIDs []ebiten.TouchID
	op       ebiten.DrawImageOptions

	Squares            []*entities.Square
	Gravity            float64
	PreviousUpdateTime time.Time
}

func Init() *Game {
	g := &Game{}
	ebiten.SetTPS(120)
	g.Squares = []*entities.Square{}
	g.Squares = append(g.Squares,
		entities.NewSquare(255, 0, 0, 90, 100, 300, 200, 0.11, -0.09, 0, 0, 1, 0.9),
		entities.NewSquare(0, 255, 0, 90, 120, 100, 500, 0.05, 0.1, 0, 0, 1, 0.9),
	)
	g.Gravity = 0.0005
	g.PreviousUpdateTime = time.Now()

	return g
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
		square.Update(timeDelta, g.Gravity)
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
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
