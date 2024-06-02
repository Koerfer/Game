package font

import (
	_ "embed"
)

var (
	//go:embed MonoBold.ttf
	MonoBold []byte

	//go:embed MonoRegular.ttf
	MonoRegular []byte
)
