package game

import (
	"encoding/gob"
	"game/game/cards"
	"game/game/screen"
	"log"
	"os"
	"time"
)

type SaveState struct {
	Cards      []*CardSaveState
	WindowSize *WindowSize
	PlayState  *PlayState
	Screen     screen.Screen
}

type CardSaveState struct {
	Id       int
	Number   int
	State    cards.State
	PlayCard *PlayCardSaveState
}

type PlayCardSaveState struct {
	Active            bool
	ActiveTime        time.Duration
	ActiveRemaining   time.Duration
	CoolDown          time.Duration
	CoolDownRemaining time.Duration

	ActiveSingleTargetDamageBoost float64
	PassiveDamageBoost            float64

	ActiveMultiTargetBoost  int
	PassiveMultiTargetBoost int
}

type PlayState struct {
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
}

func (g *Game) Save() {
	var cardsState []*CardSaveState
	var selectedCardIds []int
	for i, selectedCard := range g.Cards.Selected {
		if selectedCard == nil {
			continue
		}
		selectedCardIds = append(selectedCardIds, selectedCard.Id)
		playCardState := &PlayCardSaveState{
			Active:                        selectedCard.Active,
			ActiveTime:                    selectedCard.ActiveTime,
			ActiveRemaining:               selectedCard.ActiveRemaining,
			CoolDown:                      selectedCard.CoolDown,
			CoolDownRemaining:             selectedCard.CoolDownRemaining,
			ActiveSingleTargetDamageBoost: selectedCard.ActiveSingleTargetDamageBoost,
			PassiveDamageBoost:            selectedCard.PassiveDamageBoost,
			ActiveMultiTargetBoost:        selectedCard.ActiveMultiTargetBoost,
			PassiveMultiTargetBoost:       selectedCard.PassiveMultiTargetBoost,
		}

		cardState := &CardSaveState{
			Id:       selectedCard.Id,
			Number:   i,
			State:    cards.StateSelected,
			PlayCard: playCardState,
		}

		cardsState = append(cardsState, cardState)
	}

	for _, card := range g.Cards.Cards {
		if card == nil {
			continue
		}
		alreadyCovered := false
		for _, id := range selectedCardIds {
			if card.Id == id {
				alreadyCovered = true
				break
			}
		}
		if alreadyCovered {
			continue
		}

		playCardState := &PlayCardSaveState{
			Active:                        card.PlayCard.Active,
			ActiveTime:                    card.PlayCard.ActiveTime,
			ActiveRemaining:               card.PlayCard.ActiveRemaining,
			CoolDown:                      card.PlayCard.CoolDown,
			CoolDownRemaining:             card.PlayCard.CoolDownRemaining,
			ActiveSingleTargetDamageBoost: card.PlayCard.ActiveSingleTargetDamageBoost,
			PassiveDamageBoost:            card.PlayCard.PassiveDamageBoost,
			ActiveMultiTargetBoost:        card.PlayCard.ActiveMultiTargetBoost,
			PassiveMultiTargetBoost:       card.PlayCard.PassiveMultiTargetBoost,
		}

		cardState := &CardSaveState{
			Id:       card.Id,
			Number:   0,
			State:    card.State,
			PlayCard: playCardState,
		}

		cardsState = append(cardsState, cardState)
	}

	var playState *PlayState
	if g.PlayState != nil {
		playState = &PlayState{
			Wave:                     g.PlayState.Wave,
			TimeRemaining:            g.PlayState.TimeRemaining,
			Playing:                  g.PlayState.Playing,
			NumberOfMonsters:         g.PlayState.NumberOfMonsters,
			NumberOfMonstersExact:    g.PlayState.NumberOfMonstersExact,
			HPPerMonster:             g.PlayState.HPPerMonster,
			MonstersKilled:           g.PlayState.MonstersKilled,
			MonstersRemaining:        g.PlayState.MonstersRemaining,
			MonsterHealth:            g.PlayState.MonsterHealth,
			NumberOfMonstersAttacked: g.PlayState.NumberOfMonstersAttacked,
			DamagePerSecond:          g.PlayState.DamagePerSecond,
			SingleTargetBoost:        g.PlayState.SingleTargetBoost,
		}
	}

	saveState := &SaveState{
		Cards:      cardsState,
		WindowSize: g.WindowSize,
		PlayState:  playState,
		Screen:     g.Screen,
	}

	dumpFile, err := os.Create("game/save/save.bin")
	if err != nil {
		log.Fatalf("unable to create data.bin file: %v", err)
	}
	defer dumpFile.Close()

	enc := gob.NewEncoder(dumpFile)
	if err := enc.Encode(saveState); err != nil {
		log.Fatalf("failing to encode data: %v", err)
	}
}

func (g *Game) Load() *SaveState {
	_, err := os.Stat("game/save/save.bin")
	if os.IsNotExist(err) {
		return nil
	}

	binaryData, err := os.Open("game/save/save.bin")
	if err != nil {
		log.Fatalf("unable to open binary file: %v", err)
	}
	defer binaryData.Close()

	saveState := SaveState{}
	dec := gob.NewDecoder(binaryData)
	if err := dec.Decode(&saveState); err != nil {
		log.Fatalf("failing to decode data: %v", err)
	}

	return &saveState
}
