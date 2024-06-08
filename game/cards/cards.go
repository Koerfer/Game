package cards

import (
	"game/game/colours"
	"game/game/font"
)

type Cards struct {
	Cards          []*Card
	NumberSelected int
}

func (cs *Cards) Init() {
	boldFont := font.GetBold()
	card1 := &Card{}
	card1.Init(4, 8, 4, 0, "Test", "Does something\ncool", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	card2 := &Card{}
	card2.Init(4, 8, 8, 0, "Test", "What's up boiii\nNot that much bruh\nAnna is really cute\nHAPPI", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	card3 := &Card{}
	card3.Init(4, 8, 12, 0, "Test", "Does something\ncool", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	card4 := &Card{}
	card4.Init(4, 8, 4, 8, "Test", "Does something\ncool", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	card5 := &Card{}
	card5.Init(4, 8, 8, 8, "Test", "Does something\ncool", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	card6 := &Card{}
	card6.Init(4, 8, 12, 8, "Test", "Does something\ncool", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	cs.Cards = append(cs.Cards, card1, card2, card3, card4, card5, card6)
}
