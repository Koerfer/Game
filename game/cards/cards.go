package cards

import (
	"game/game/colours"
	"game/game/font"
	"time"
)

type Cards struct {
	Cards          []*Card
	Selected       []*Card
	NumberSelected int
}

func (cs *Cards) Init() {
	boldFont := font.GetBold()
	card1 := &Card{}
	card1.Init(4, 0, "Damage", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	card2 := &Card{}
	card2.Init(8, 0, "Multi Target", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	card3 := &Card{}
	card3.Init(12, 0, "Chronos", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	card4 := &Card{}
	card4.Init(4, 8, "Cutie", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	card5 := &Card{}
	card5.Init(8, 8, "Test", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	card6 := &Card{}
	card6.Init(12, 8, "Test", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)

	card1.addEffect("Passive: Increases damage by %dx\nActive: 10s of %dx single\ntarget damage\nCooldown: %02d s", 2, 5, 10*time.Second, 30*time.Second)
	card2.addEffect("Hello cool girl", 1, 1, 10*time.Second, 30*time.Second)
	card3.addEffect("Testing TESTING TeStInG", 1, 1, 10*time.Second, 30*time.Second)
	card4.addEffect("Anna Cutie is the\nBEST", 1, 1, 10*time.Second, 30*time.Second)
	card5.addEffect("I should really go to sleep", 1, 1, 10*time.Second, 30*time.Second)
	card6.addEffect("Game is getting there", 1, 1, 10*time.Second, 30*time.Second)
	cs.Cards = append(cs.Cards, card1, card2, card3, card4, card5, card6)
	cs.Selected = make([]*Card, 3)
}
