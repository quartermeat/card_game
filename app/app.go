// Package 'app' is the main game function and main loop
package app

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
	"github.com/quartermeat/card_game/debuglog"
	"github.com/quartermeat/card_game/input"
	"github.com/quartermeat/card_game/objects"
)

// App holds the main game loop
func App() {

	cfg := pixelgl.WindowConfig{
		Title:  "Card Game",
		Bounds: pixel.R(0, 0, 1290, 1080),
		VSync:  false,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var (
		camPos             = pixel.ZV
		camSpeed           = 500.0
		camZoom            = 1.0
		camZoomSpeed       = 1.2
		gameObjs           objects.GameObjects
		gameCommands       = make(input.Commands)
		frames             = 0
		second             = time.NewTicker(time.Second)
		drawHitBox         = false
		inputHandler       input.InputHandler
		objectAssets       assets.ObjectAssets
		debugLog           debuglog.Entries
		sysErrors          []error
		consoleToInputChan chan console.IConsoleTxTopic
	)

	consoleToInputChan = make(chan console.IConsoleTxTopic, 1)
	defer close(consoleToInputChan)

	// start command server
	go console.StartServer(consoleToInputChan)
	if Test {
		go console.RunConsole()
	}

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

		debugLog = inputHandler.HandleInput(
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
			consoleToInputChan,
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
		//must be done outside of inputHandler to be the last thing drawn
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
		case <-second.C:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | GameObjects: %d", cfg.Title, frames, len(gameObjs)))
			frames = 0
		default:
		}

		//output errors every loop
		for _, entry := range debugLog {
			fmt.Printf("debugLog: %s", entry.GetMessage())
			//gracefully handle the stop command
			if entry.GetMessage() == console.Stop {
				win.Destroy()
			}
		}
	}
}
