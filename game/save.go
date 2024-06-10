package game

import (
	"encoding/gob"
	"game/game/cards"
	"game/game/play"
	"game/game/screen"
	"log"
	"os"
	"time"
)

type SaveState struct {
	Cards           []*CardSaveState
	MenuItemShown   map[int]bool
	WindowSize      *WindowSize
	PlayState       *PlayState
	Screen          screen.Screen
	StartButtonName string
	HighestWave     int
	Upgrades        int
	NumberSelected  int
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

	PassiveTimeSlow int64
	ActiveTimeSkip  time.Duration
}

type PlayState struct {
	Wave          int
	TimeRemaining time.Duration
	Playing       bool

	NumberOfMonsters      int
	NumberOfMonstersExact float64
	HPPerMonster          float64
	HPPerMonsterCopy      float64

	MonstersKilled    int
	MonstersRemaining int
	MonsterHealth     []float64

	NumberOfMonstersAttacked int
	DamagePerSecond          float64
	SingleTargetBoost        float64
	TimeSlow                 int64
	TimeSkip                 time.Duration
}

func (g *Game) Save() {
	var cardsState []*CardSaveState
	var selectedCardIds []int
	for i, selectedCard := range g.Cards.Selected {
		if selectedCard == nil {
			continue
		}
		selectedCardIds = append(selectedCardIds, selectedCard.Id)
		playCardState := playCardStateConvert(selectedCard)

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

		playCardState := playCardStateConvert(card.PlayCard)

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
			HPPerMonsterCopy:         g.PlayState.HPPerMonsterCopy,
			MonstersKilled:           g.PlayState.MonstersKilled,
			MonstersRemaining:        g.PlayState.MonstersRemaining,
			MonsterHealth:            g.PlayState.MonsterHealth,
			NumberOfMonstersAttacked: g.PlayState.NumberOfMonstersAttacked,
			DamagePerSecond:          g.PlayState.DamagePerSecond,
			SingleTargetBoost:        g.PlayState.SingleTargetBoost,
			TimeSlow:                 g.PlayState.TimeSlow,
			TimeSkip:                 g.PlayState.TimeSkip,
		}
	}

	saveState := &SaveState{
		Cards:           cardsState,
		WindowSize:      g.WindowSize,
		PlayState:       playState,
		Screen:          g.Screen,
		StartButtonName: g.StartButton.Name,
		HighestWave:     g.HighestWave,
		MenuItemShown:   make(map[int]bool),
		Upgrades:        g.Cards.Upgrades,
		NumberSelected:  g.Cards.NumberSelected,
	}

	for i, menuItem := range g.MenuItems {
		saveState.MenuItemShown[i] = menuItem.Shown
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

func playCardStateConvert(card *cards.PlayCard) *PlayCardSaveState {
	return &PlayCardSaveState{
		Active:                        card.Active,
		ActiveTime:                    card.ActiveTime,
		ActiveRemaining:               card.ActiveRemaining,
		CoolDown:                      card.CoolDown,
		CoolDownRemaining:             card.CoolDownRemaining,
		ActiveSingleTargetDamageBoost: card.ActiveSingleTargetDamageBoost,
		PassiveDamageBoost:            card.PassiveDamageBoost,
		ActiveMultiTargetBoost:        card.ActiveMultiTargetBoost,
		PassiveMultiTargetBoost:       card.PassiveMultiTargetBoost,
		PassiveTimeSlow:               card.PassiveTimeSlow,
		ActiveTimeSkip:                card.ActiveTimeSkip,
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

func (g *Game) UpdateToMatchLoadedState(state *SaveState) {
	g.WindowSize = state.WindowSize
	g.Screen = state.Screen
	g.StartButton.Name = state.StartButtonName
	g.HighestWave = state.HighestWave

	for i, menuItem := range g.MenuItems {
		menuItem.Shown = state.MenuItemShown[i]
	}

	if state.PlayState != nil {
		g.PlayState = &play.State{}
		g.PlayState.Wave = state.PlayState.Wave
		g.PlayState.TimeRemaining = state.PlayState.TimeRemaining
		g.PlayState.Playing = state.PlayState.Playing
		g.PlayState.NumberOfMonsters = state.PlayState.NumberOfMonsters
		g.PlayState.NumberOfMonstersExact = state.PlayState.NumberOfMonstersExact
		g.PlayState.HPPerMonster = state.PlayState.HPPerMonster
		g.PlayState.HPPerMonsterCopy = state.PlayState.HPPerMonsterCopy
		g.PlayState.MonstersKilled = state.PlayState.MonstersKilled
		g.PlayState.MonstersRemaining = state.PlayState.MonstersRemaining
		g.PlayState.MonsterHealth = state.PlayState.MonsterHealth
		g.PlayState.NumberOfMonstersAttacked = state.PlayState.NumberOfMonstersAttacked
		g.PlayState.DamagePerSecond = state.PlayState.DamagePerSecond
		g.PlayState.SingleTargetBoost = state.PlayState.SingleTargetBoost
		g.PlayState.TimeSlow = state.PlayState.TimeSlow
		g.PlayState.TimeSkip = state.PlayState.TimeSkip
		g.PlayState.ActiveCards = make([]*cards.PlayCard, 3)
	}

	g.Cards.Upgrades = state.Upgrades
	g.Cards.NumberSelected = state.NumberSelected
	for _, card := range g.Cards.Cards {
		for _, saveCard := range state.Cards {
			if card.Id != saveCard.Id {
				continue
			}

			card.State = saveCard.State
			card.PlayCard.Active = saveCard.PlayCard.Active
			card.PlayCard.ActiveTime = saveCard.PlayCard.ActiveTime
			card.PlayCard.ActiveRemaining = saveCard.PlayCard.ActiveRemaining
			card.PlayCard.CoolDown = saveCard.PlayCard.CoolDown
			card.PlayCard.CoolDownRemaining = saveCard.PlayCard.CoolDownRemaining
			card.PlayCard.ActiveSingleTargetDamageBoost = saveCard.PlayCard.ActiveSingleTargetDamageBoost
			card.PlayCard.PassiveDamageBoost = saveCard.PlayCard.PassiveDamageBoost
			card.PlayCard.ActiveMultiTargetBoost = saveCard.PlayCard.ActiveMultiTargetBoost
			card.PlayCard.PassiveMultiTargetBoost = saveCard.PlayCard.PassiveMultiTargetBoost
			card.PlayCard.PassiveTimeSlow = saveCard.PlayCard.PassiveTimeSlow
			card.PlayCard.ActiveTimeSkip = saveCard.PlayCard.ActiveTimeSkip
			card.Update(g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)

			if card.State == cards.StateSelected {
				card.AddToHand(saveCard.Number+1, g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)
				g.Cards.Selected[saveCard.Number] = card.PlayCard
			}

			if g.PlayState != nil && card.PlayCard.Active {
				g.PlayState.ActiveCards[saveCard.Number] = card.PlayCard
			}
		}
	}

	for _, menuItem := range g.MenuItems {
		menuItem.UpdateSize(g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)
	}
	for _, card := range g.Cards.Cards {
		card.Update(g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)
	}
	for _, divider := range g.MainDividers {
		divider.UpdateSize(g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)
	}
	for _, divider := range g.PlayDividers {
		divider.UpdateSize(g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)
	}
	g.StartButton.UpdateSize(g.WindowSize.CurrentWidthFactor, g.WindowSize.CurrentHeightFactor)

	g.WindowSize.PreviousHeightFactor = g.WindowSize.CurrentHeightFactor
	g.WindowSize.PreviousWidthFactor = g.WindowSize.CurrentWidthFactor
}
