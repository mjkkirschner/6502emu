package main

import (
	"image/color"
)

//method to iterate screen memory in character mode

func (sim *Simulator) GetColorDataForScreenInCharacterMode() []color.RGBA {
	//there are 1000 bytes of screen memory, but each character is an 8*8 pixel memory region (64 pixels long.)
	out := make([]color.RGBA, 64000)
	for i := range out {
		out[i] = color.RGBA{255, 255, 255, 255}
	}
	////
	currentScreenX := 0
	currentScreenY := 0
	for currentScreenPos := 0; currentScreenPos < 1000; currentScreenPos++ {
		currentCharCode := sim.Memory[memoryMap["SCREEN"].start+uint16(currentScreenPos)]

		if currentScreenX == 40 {
			currentScreenX = 0
			currentScreenY++
		}
		//8 rows per character in char rom map.
		for currentRow := 0; currentRow < 8; currentRow++ {
			currentLine := sim.Memory[memoryMap["CHAR"].start+(uint16(currentCharCode)<<3)+uint16(currentRow)]
			//8 cols per character in char rom map
			for currentCharacterCol := 0; currentCharacterCol < 8; currentCharacterCol++ {
				//TODO why 128 is on... because there is padding around the character?
				isPixelOn := (currentLine & 0x80) == 0x80
				pixelx := (currentScreenX << 3) + currentCharacterCol
				pixely := (currentScreenY << 3) + currentRow
				bufferIndex := (pixely*320 + pixelx)
				if isPixelOn {
					out[bufferIndex] = color.RGBA{0, 0, 0, 255}

				} else {
					out[bufferIndex] = color.RGBA{255, 255, 255, 255}
				}
				currentLine = currentLine << 1

			}
		}
		currentScreenX++
	}
	return out
}
