package helper

import "math"

func GetNewTextSize(baseTextSize, heightFactor, newWidth float64, text string) float64 {
	newNameTextSize := baseTextSize * heightFactor
	if newNameTextSize*float64(len(text)) > newWidth {
		newNameTextSize = math.Min(newNameTextSize, 1.4*newWidth/float64(len(text)))
	}

	return newNameTextSize
}
