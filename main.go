package main

import (
	"fmt"
	_ "image/png"
	"math/rand"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/quartermeat/aeonExMachina/assets"
	"github.com/quartermeat/aeonExMachina/input"
	"github.com/quartermeat/aeonExMachina/objects"
	"golang.org/x/image/colornames"
)

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "Aeon Ex Machina",
		Bounds: pixel.R(0, 0, 1280, 960),
		VSync:  false,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// start command server
	go StartServer()

	var (
		camPos             = pixel.ZV
		camSpeed           = 500.0
		camZoom            = 1.0
		camZoomSpeed       = 1.2
		gameObjs           objects.GameObjects
		gameCommands       = make(input.Commands)
		frames             = 0
		second             = time.Tick(time.Second)
		drawHitBox         = false
		inputHandler       input.InputHandler
		livingObjectAssets assets.ObjectAssets
		gibletObjectAssets assets.ObjectAssets
	)

	//load assets
	err = livingObjectAssets.SetAssets("assets/spriteSheet.png", "assets/pinkAnimations.csv", 32)
	if err != nil {
		panic(err)
	}
	err = gibletObjectAssets.SetAssets("assets/spriteSheet.png", "assets/gibletAnimations.csv", 16)
	if err != nil {
		panic(err)
	}

	//seed rng
	rand.Seed(time.Now().UnixNano())

	last := time.Now()
	for !win.Closed() {
		//handle delta
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		inputHandler.HandleInput(
			win,
			&cam,
			gameCommands,
			&gameObjs,
			gibletObjectAssets,
			livingObjectAssets,
			dt,
			camSpeed,
			&camZoom,
			camZoomSpeed,
			&camPos,
			&drawHitBox,
		)

		var waitGroup sync.WaitGroup

		//handle game updates
		gameCommands.ExecuteCommands(&waitGroup)
		waitGroup.Wait()
		gameObjs.UpdateAllObjects(dt, &waitGroup)
		waitGroup.Wait()

		win.Clear(colornames.Black)
		//draw game objects
		gameObjs.DrawAllObjects(win, drawHitBox, &waitGroup)
		waitGroup.Wait()

		//draw cursor based on selected object
		if win.MouseInsideWindow() {
			if !win.Pressed(pixelgl.KeyLeftControl) {
				win.SetCursorVisible(false)
				inputHandler.ObjectToPlace.Sprite().Draw(win, pixel.IM.Moved(cam.Unproject(win.MousePosition())))
			}
		} else {
			win.SetCursorVisible(true)
		}

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
