package font

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

var (
	//go:embed MonoBold.ttf
	monoBold []byte

	//go:embed MonoRegular.ttf
	monoRegular []byte
)

func GetBold() *text.GoTextFaceSource {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(monoBold))
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func GetRegular() *text.GoTextFaceSource {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(monoRegular))
	if err != nil {
		log.Fatal(err)
	}
	return s
}
