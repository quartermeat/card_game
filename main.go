package main

import (
	"fmt"
	_ "image/png"
	"math/rand"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/quartermeat/card_game/assets"
	"github.com/quartermeat/card_game/console"
	"github.com/quartermeat/card_game/errormgmt"
	"github.com/quartermeat/card_game/input"
	"github.com/quartermeat/card_game/objects"
)

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "Card Game",
		Bounds: pixel.R(0, 0, 1290, 1080),
		VSync:  false,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// start command server
	go console.StartServer()

	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
		gameObjs     objects.GameObjects
		gameCommands = make(input.Commands)
		frames       = 0
		second       = time.Tick(time.Second)
		drawHitBox   = false
		inputHandler input.InputHandler
		objectAssets assets.ObjectAssets
		errors       errormgmt.Errors
		sysErrors    []error
	)

	//panic level errors
	sysErrors = make([]error, 0)

	//load assets
	objectAssets, err1 := objectAssets.AddAssets(assets.CursorAnimations, "assets/mouseHand.png", "assets/mouseAnimations.csv", assets.MouseIconPixelSize)
	objectAssets, err2 := objectAssets.AddAssets(assets.TestCard, "assets/test_card.png", "assets/testCardAnimations.csv", assets.CardImageSize)
	sysErrors = append(sysErrors, err1)
	sysErrors = append(sysErrors, err2)
	for _, sysError := range sysErrors {
		if sysError != nil {
			panic(sysError)
		}
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

		errors = inputHandler.HandleInput(
			win,
			&cam,
			gameCommands,
			&gameObjs,
			objectAssets,
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
				//setup and object to place
				inputHandler.Cursor.Draw(win, pixel.IM.Moved(cam.Unproject(win.MousePosition())))
			}
		} else {
			win.SetCursorVisible(true)
		}

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | GameObjects: %d", cfg.Title, frames, len(gameObjs)))
			frames = 0
		default:
		}

		//output errors every loop
		//TODO: set output to a debug window
		for i, error := range errors {
			fmt.Printf("error %d: %s", i, error.Error())
		}
	}
}

func main() {
	// fmt.Printf("hello go 1.19")
	pixelgl.Run(run)
}
