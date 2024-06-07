package helper

func GetNewSizeAndPosition(baseWidth, baseHeight, basePosX, basePosY, widthFactor, heightFactor, shiftX, shiftY float64) (float64, float64, float64, float64) {
	newWidth := baseWidth * widthFactor
	newHeight := baseHeight * heightFactor
	newX := basePosX*widthFactor + shiftX
	newY := basePosY*heightFactor + shiftY

	return newWidth, newHeight, newX, newY
}
