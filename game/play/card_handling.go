package play

import "game/game/cards"

func (s *State) CardActivation(card *cards.Card, number int) {
	card.PlayCard.Active = true
	card.PlayCard.ActiveRemaining = card.ActivationTime
	s.SingleTargetBoost *= card.ActiveSingleTargetDamageBoost
	s.ActiveCards[number] = card
}

func (s *State) deactivateCard(card *cards.Card) {
	card.PlayCard.Active = false
	card.PlayCard.ActiveRemaining = 0
	card.PlayCard.CoolDownRemaining = card.CoolDown
	s.SingleTargetBoost /= card.ActiveSingleTargetDamageBoost

}
