package main

import (
	_ "image/png"

	"github.com/faiface/pixel/pixelgl"
	"github.com/quartermeat/card_game/app"
)

func main() {
	pixelgl.Run(app.App)
	//scratch.Run()
}
