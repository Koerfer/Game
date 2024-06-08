package screen

type Screen uint8

const (
	ScreenInvalid Screen = iota
	ScreenMain
	ScreenCards
	ScreenTech
	ScreenAnna
	ScreenSettings
)
