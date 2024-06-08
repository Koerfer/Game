package screen

type Screen uint8

const (
	ScreenNothing Screen = iota
	ScreenMain
	ScreenCards
	ScreenTech
	ScreenAnna
	ScreenSettings
	ScreenPlay
)
