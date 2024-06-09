package play

import (
	"game/game/cards"
	"time"
)

func Start(selectedCards []*cards.PlayCard) *State {
	initialState := &State{
		Wave:                     1,
		TimeRemaining:            1 * time.Minute,
		NumberOfMonsters:         1,
		Playing:                  true,
		NumberOfMonstersExact:    1,
		HPPerMonster:             2,
		MonstersRemaining:        1,
		MonsterHealth:            []float64{2},
		DamagePerSecond:          5, // 0.2
		NumberOfMonstersAttacked: 2, // 1
		ActiveCards:              make([]*cards.PlayCard, 3),
		SingleTargetBoost:        1,
	}

	for _, selectedCard := range selectedCards {
		if selectedCard == nil {
			continue
		}

		initialState.DamagePerSecond *= selectedCard.PassiveDamageBoost
		initialState.NumberOfMonstersAttacked *= selectedCard.PassiveMultiTargetBoost
	}

	return initialState
}
