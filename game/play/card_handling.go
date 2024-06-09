package play

import "game/game/cards"

func (s *State) CardActivation(card *cards.PlayCard, number int) {
	card.Active = true
	card.ActiveRemaining = card.ActiveTime
	s.SingleTargetBoost *= card.ActiveSingleTargetDamageBoost
	s.NumberOfMonstersAttacked *= card.ActiveMultiTargetBoost
	s.ActiveCards[number] = card
}

func (s *State) deactivateCard(card *cards.PlayCard) {
	card.Active = false
	card.ActiveRemaining = 0
	card.CoolDownRemaining = card.CoolDown
	s.SingleTargetBoost /= card.ActiveSingleTargetDamageBoost
	s.NumberOfMonstersAttacked /= card.ActiveMultiTargetBoost

}
