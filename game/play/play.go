package play

import (
	"game/game/cards"
	"time"
)

type State struct {
	Wave          int
	TimeRemaining time.Duration
	Playing       bool

	NumberOfMonsters      int
	NumberOfMonstersExact float64
	HPPerMonster          float64

	MonstersKilled    int
	MonstersRemaining int
	MonsterHealth     []float64

	NumberOfMonstersAttacked int
	DamagePerSecond          float64
	SingleTargetBoost        float64

	ActiveCards []*cards.PlayCard
}

func (s *State) Update(timeDelta time.Duration) {
	if s.TimeRemaining <= 0 {
		s.Playing = false
		s.TimeRemaining = 0
		return
	}

	s.TimeRemaining -= timeDelta

	for i, activeCard := range s.ActiveCards {
		if activeCard == nil {
			continue
		}

		if activeCard.ActiveRemaining <= 0 {
			s.deactivateCard(activeCard)
			s.ActiveCards[i] = nil
		}
	}

	if s.MonstersRemaining == 0 {
		s.prepareNewWave()
	}

	killed := s.MonstersKilled
	upTo := s.NumberOfMonstersAttacked + killed
	if s.NumberOfMonstersAttacked+killed > s.NumberOfMonsters {
		upTo = s.NumberOfMonsters
	}

	for i := killed; i < upTo; i++ {
		if i == killed {
			s.MonsterHealth[i] -= s.DamagePerSecond * timeDelta.Seconds() * s.SingleTargetBoost
		} else {
			s.MonsterHealth[i] -= s.DamagePerSecond * timeDelta.Seconds()
		}

		if s.MonsterHealth[i] <= 0 {
			s.MonstersRemaining -= 1
			s.MonstersKilled += 1
		}
	}
}

func (s *State) prepareNewWave() {
	s.Wave++
	s.MonstersKilled = 0
	s.NumberOfMonstersExact *= 1.15
	s.NumberOfMonsters = int(s.NumberOfMonstersExact)
	s.HPPerMonster *= 1.2
	s.MonstersRemaining = s.NumberOfMonsters
	healths := make([]float64, s.NumberOfMonsters)
	for i := range healths {
		healths[i] = s.HPPerMonster
	}
	s.MonsterHealth = healths
}
