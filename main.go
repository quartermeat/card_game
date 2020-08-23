package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "Aeon Ex Machina",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  false,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	go StartServer()

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
			if win.Pressed(pixelgl.KeyLeftControl) {
				mouse := cam.Unproject(win.MousePosition())
				selectedObj, err := gameObjs.getSelectedGameObj(mouse)
				if err != nil {
					fmt.Printf(err.Error())
				}
				fmt.Println("object id:", selectedObj.id)
			} else {
				mouse := cam.Unproject(win.MousePosition())
				gameObjs = gameObjs.addGameObject(pinkAnimKeys, pinkAnims, pinkSheet, mouse)
			}
		}

		if win.Pressed(pixelgl.MouseButtonLeft) {
			if win.Pressed(pixelgl.KeyLeftShift) {
				mouse := cam.Unproject(win.MousePosition())
				gameObjs = gameObjs.addGameObject(pinkAnimKeys, pinkAnims, pinkSheet, mouse)
			}
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

		win.Clear(colornames.Black)

		//this is craziness
		var waitGroup sync.WaitGroup

		//handle updates
		gameObjs.updateAll(dt, &waitGroup)
		waitGroup.Wait()
		//handle drawing
		gameObjs.drawAll(win, &waitGroup)
		waitGroup.Wait()

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
