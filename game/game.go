package game

import (
	"bytes"
	"fmt"
	"game/font"
	"image/color"
	"log"
	"math/rand"
	"time"

	"game/game/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Game struct {
	touchIDs []ebiten.TouchID
	op       ebiten.DrawImageOptions
	font     *text.GoTextFaceSource

	Squares            []*entities.Square
	Gravity            float64
	PreviousUpdateTime time.Time

	LastBounceTime  time.Time
	LastAddedSquare time.Time
	HighScore       time.Duration
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

	s, err := text.NewGoTextFaceSource(bytes.NewReader(font.MonoBold))
	if err != nil {
		log.Fatal(err)
	}
	g.font = s

	g.LastBounceTime = time.Now()
	g.LastAddedSquare = time.Now()

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

	timeSinceLastBounce := time.Since(g.LastBounceTime)
	var removeAdditionalSquares bool
	for _, square := range g.Squares {
		bounced := square.Update(timeDelta, g.Gravity)
		if bounced {
			removeAdditionalSquares = true
			if g.HighScore < timeSinceLastBounce {
				g.HighScore = timeSinceLastBounce
			}
			g.LastBounceTime = time.Now()
			g.LastAddedSquare = time.Now()
		}
	}
	if removeAdditionalSquares {
		g.Squares = append(g.Squares[:2], g.Squares[len(g.Squares):]...)
	}

	if time.Now().Sub(g.LastAddedSquare) > time.Second*10 {
		screenWidth, screenHeight := ebiten.WindowSize()
		posX := float64(screenWidth)*0.1 + rand.Float64()*float64(screenWidth)*0.9
		posY := float64(screenHeight)*0.1 + rand.Float64()*float64(screenHeight)*0.4
		movX := (rand.Float64() - 0.5) * 0.2
		movY := (rand.Float64() - 0.5) * 0.2
		g.Squares = append(g.Squares,
			entities.NewSquare(0, 0, 255, 90, 100, posX, posY, movX, movY, 0, 0, 1, 0.9),
		)
		g.LastAddedSquare = time.Now()
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	screenWidth, _ := ebiten.WindowSize()
	for _, square := range g.Squares {
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(square.PosX, square.PosY)

		screen.DrawImage(square.Image, &ebiten.DrawImageOptions{
			GeoM: g.op.GeoM,
		})
	}

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)

	msg := fmt.Sprintf(`TPS: %0.2f - FPS: %0.2f`, ebiten.ActualTPS(), ebiten.ActualFPS())
	op.GeoM.Translate(float64(screenWidth)-170, 0)
	text.Draw(screen, msg, &text.GoTextFace{
		Source: g.font,
		Size:   12,
	}, op)
	op.GeoM.Reset()

	timeSinceBounce := time.Since(g.LastBounceTime)
	msg2 := fmt.Sprintf(`Current time: %02d:%02d:%03d`, int(timeSinceBounce.Minutes()), int(timeSinceBounce.Seconds())%60, int(timeSinceBounce.Milliseconds())%1000)
	op.GeoM.Translate(0, 0)
	text.Draw(screen, msg2, &text.GoTextFace{
		Source: g.font,
		Size:   24,
	}, op)

	msg3 := fmt.Sprintf(`Best time:     %02d:%02d:%03d`, int(g.HighScore.Minutes()), int(g.HighScore.Seconds())%60, int(g.HighScore.Milliseconds())%1000)
	op.GeoM.Translate(0, 28)
	text.Draw(screen, msg3, &text.GoTextFace{
		Source: g.font,
		Size:   24,
	}, op)
}
