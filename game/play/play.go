package play

import (
	"game/game/cards"
	"math/rand"
	"time"
)

type State struct {
	Wave          int
	TimeRemaining time.Duration
	Playing       bool

	NumberOfMonsters      int
	NumberOfMonstersExact float64
	HPPerMonster          float64
	HPPerMonsterCopy      float64 // to be returned to after boss fight

	MonstersKilled    int
	MonstersRemaining int
	MonsterHealth     []float64

	NumberOfMonstersAttacked int
	DamagePerSecond          float64
	SingleTargetBoost        float64
	TimeSlow                 int64
	TimeSkip                 time.Duration

	ActiveCards []*cards.PlayCard
}

func (s *State) Update(timeDelta time.Duration, highestWave int) int8 {
	s.TimeRemaining -= time.Duration(timeDelta.Milliseconds()/s.TimeSlow) * time.Millisecond
	if s.TimeRemaining <= -5*time.Second {
		return -1
	}
	if s.TimeRemaining <= 0 {
		s.Playing = false
		return 0
	}

	for i, activeCard := range s.ActiveCards {
		if activeCard == nil {
			continue
		}

		activeCard.ActiveRemaining -= timeDelta
		if activeCard.Active && activeCard.ActiveRemaining <= 0 {
			s.startCooldownForCard(activeCard)
			continue
		}

		if !activeCard.Active {
			activeCard.CoolDownRemaining -= timeDelta
			if activeCard.CoolDownRemaining <= 0 {
				s.ActiveCards[i] = nil
			}
		}
	}

	if s.MonstersRemaining == 0 {
		mod10 := s.Wave % 10
		s.prepareNewWave(highestWave)

		if s.Wave > highestWave {
			if mod10 == 5 {
				if s.Wave > 55 {
					return 0
					// no new cards to gain
				}
				return 1
			}
			if mod10 == 0 {
				return 2
			}
		}

		if mod10 == 0 && rand.Float64() < float64(s.Wave)/float64(highestWave) {
			return 2
		}
	}

	killed := s.MonstersKilled
	upTo := s.NumberOfMonstersAttacked + killed
	if s.NumberOfMonstersAttacked+killed > s.NumberOfMonsters {
		upTo = s.NumberOfMonsters
	}

	for i := killed; i < upTo; i++ {
		if i == killed {
			s.MonsterHealth[i] -= s.DamagePerSecond * (timeDelta.Seconds() + s.TimeSkip.Seconds()) * s.SingleTargetBoost
		} else {
			s.MonsterHealth[i] -= s.DamagePerSecond * (timeDelta.Seconds() + s.TimeSkip.Seconds())
		}
		s.TimeSkip = 0

		if s.MonsterHealth[i] <= 0 {
			s.MonstersRemaining -= 1
			s.MonstersKilled += 1
		}
	}

	return 0
}

func (s *State) prepareNewWave(highestWave int) {
	if s.Wave == highestWave {
		s.TimeSlow = 1
	}
	s.Wave++
	mod10 := s.Wave % 10

	s.MonstersKilled = 0
	s.NumberOfMonstersExact *= 1.08
	s.HPPerMonsterCopy *= 1.12

	switch mod10 {
	case 0, 5:
		s.NumberOfMonsters = 1
		s.HPPerMonster = s.HPPerMonsterCopy * 5
		s.HPPerMonsterCopy *= 1.12
	case 1, 6:
		s.HPPerMonster = s.HPPerMonsterCopy
		s.NumberOfMonsters = int(s.NumberOfMonstersExact)
	default:
		s.NumberOfMonsters = int(s.NumberOfMonstersExact)
	}

	s.HPPerMonster *= 1.12
	s.MonstersRemaining = s.NumberOfMonsters
	healths := make([]float64, s.NumberOfMonsters)
	for i := range healths {
		healths[i] = s.HPPerMonster
	}
	s.MonsterHealth = healths
}
