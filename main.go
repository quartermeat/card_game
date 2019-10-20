package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func imd_draw(image *imdraw.IMDraw) {
	image.Color = colornames.Red
	image.EndShape = imdraw.RoundEndShape
	image.Push(pixel.V(100, 100), pixel.V(700, 100))
	image.Line(30)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "PixelLifeGo",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	for !win.Closed() {
		imd.Clear()
		imd_draw(imd)
		win.Clear(colornames.Mediumspringgreen)
		imd.Draw(win)
		win.Update()
	}
}

func main() {
	go func() {
		fmt.Printf(welcome_string)
	}()
	pixelgl.Run(run)
}
