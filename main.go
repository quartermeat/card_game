package main

import (
	"fmt"
	"math"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Life in Go",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	//game variables
	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
		gameObjs     GameObjects
		frames       = 0
		second       = time.Tick(time.Second)
	)

	//load assets
	pinkSheet, pinkAnims, pinkAnimKeys, err := loadAnimationSheet("assets/pink.png", "assets/pink_animations.csv", 32)
	if err != nil {
		panic(err)
	}

	last := time.Now()
	for !win.Closed() {
		//handle delta
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		//handle input
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			mouse := cam.Unproject(win.MousePosition())
			gameObjs = gameObjs.addGameObject(pinkAnimKeys, pinkAnims, pinkSheet, pixel.IM.Scaled(pixel.ZV, 1).Moved(mouse))
		}

		if win.Pressed(pixelgl.KeyA) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyD) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyS) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyW) {
			camPos.Y += camSpeed * dt
		}
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		//handle updates
		gameObjs.updateAll(dt)

		win.Clear(colornames.Black)

		//handle drawing
		gameObjs.drawAll(win)

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | CELLS: %d", cfg.Title, frames, len(gameObjs)))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
