package cards

import (
	"game/game/colours"
	"game/game/font"
	"time"
)

type Cards struct {
	Cards          []*Card
	Selected       []*PlayCard
	NumberSelected int

	Upgrades int
}

func (cs *Cards) Init() {
	boldFont := font.GetBold()
	card1 := &Card{}
	card1.Init(1, 4, 0, "Damage", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	card2 := &Card{}
	card2.Init(2, 8, 0, "Multi Target", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	card3 := &Card{}
	card3.Init(3, 12, 0, "Chronos", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	card4 := &Card{}
	card4.Init(4, 4, 8, "Test", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	card5 := &Card{}
	card5.Init(5, 8, 8, "Test", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)
	card6 := &Card{}
	card6.Init(6, 12, 8, "Test", boldFont, 2, 1, colours.White, colours.Black, colours.Blue)

	card1.PlayCard.addEffect("Passive: Increases damage by %.3g\nActive: 10s of %.3g times\nhigher single target damage\nCooldown: %02d s", 2, 4,
		1, 1,
		1, 0, 10*time.Second, 30*time.Second)
	card2.PlayCard.addEffect("Passive: Increases number of\ntargets by %.3g\nActive: Increases number of\ntargets by a further %dx\nfor 5s\nCooldown: %02d s", 1, 1,
		2, 2,
		1, 0, 5*time.Second, 30*time.Second)
	card3.PlayCard.addEffect("Passive: Time remaining goes %d times\nslower until highest wave reached\nActive: Deal %d seconds of\ndamage immediately\nCooldown: %02d s", 1, 1,
		1, 1,
		2, 4*time.Second, 0, 10*time.Second)
	card4.PlayCard.addEffect("TBD", 1, 1,
		1, 1,
		1, 0, 10*time.Second, 30*time.Second)
	card5.PlayCard.addEffect("TBD", 1, 1,
		1, 1,
		1, 0, 10*time.Second, 30*time.Second)
	card6.PlayCard.addEffect("TBD", 1, 1,
		1, 1,
		1, 0, 10*time.Second, 30*time.Second)
	cs.Cards = append(cs.Cards, card1, card2, card3, card4, card5, card6)
	cs.Selected = make([]*PlayCard, 3)
}
