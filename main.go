package main

import (
	_ "image/png"

	"github.com/faiface/pixel/pixelgl"
	"github.com/quartermeat/card_game/app"
)

// utilizes the Pixel library for 2D game development and a custom package for the card game logic. The main function calls the pixelgl.Run function with app.AppRun as an argument, which will run the card game in an OpenGL-backed window with input handling.
func main() {
	pixelgl.Run(app.AppRun)
	// scratch.RunDalleTest()
}
