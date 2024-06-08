package play

import "time"

func Start() *State {
	initialState := &State{
		Wave:              1,
		TimeRemaining:     1 * time.Minute,
		NumberOfMonsters:  1,
		HPPerMonster:      5,
		MonstersRemaining: 5,
		MonsterHealth:     []int{5},
	}

	return initialState
}
