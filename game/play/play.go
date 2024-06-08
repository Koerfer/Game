package play

import "time"

type State struct {
	Wave          int
	TimeRemaining time.Duration

	NumberOfMonsters int
	HPPerMonster     int

	MonstersRemaining int
	MonsterHealth     []int
}
