// Package 'app' is the main game function and main loop
package app

import (
	_ "image/png"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/quartermeat/card_game/app/domainApp/rules"
	"github.com/quartermeat/card_game/assets"
)

// AppRun() is the main game function and main loop for a card game.
// It sets up the window configuration, initializes the GUI, loads assets,
// seeds the random number generator, and starts a command server.
// It then enters a loop that handles delta time, handles input, updates game objects,
// draws game objects, draws the GUI, and draws a cursor based on selected object.
// At the end of each loop it also updates the window title with FPS and number of game objects.
// Finally it checks for any debug log entries with a message of 'console.Stop' and closes the window if found.
func AppRun() {

	//seed rng
	rand.New(rand.NewSource(time.Now().UnixNano()))

	cfg := pixelgl.WindowConfig{
		Title:  rules.APP_TITLE,
		Bounds: rules.WINDOW_SIZE,
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Replace the path with the path to your wooden texture image
	woodenTexture, err := assets.LoadPicture(rules.BACKGROUND_IMAGE)
	if err != nil {
		panic(err)
	}

	woodenSprite := pixel.NewSprite(woodenTexture, woodenTexture.Bounds())

	for !win.Closed() {
		win.Clear(colornames.Skyblue)

		// Draw the wooden background
		scaleX := win.Bounds().W() / woodenTexture.Bounds().W()
		scaleY := win.Bounds().H() / woodenTexture.Bounds().H()
		mat := pixel.IM.Scaled(pixel.ZV, scaleX).ScaledXY(pixel.ZV, pixel.V(scaleX, scaleY))
		woodenSprite.Draw(win, mat)
		// Add your drawing code here.

		win.Update()
	}

	// var (
	// 	appState           = observable.ObservableState{}
	// 	camPos             = pixel.ZV
	// 	camSpeed           = 500.0
	// 	camZoom            = 1.0
	// 	camZoomSpeed       = 1.2
	// 	gameObjs           objects.GameObjects
	// 	gameCommands       = make(input.Commands)
	// 	frames             = 0
	// 	second             = time.NewTicker(time.Second)
	// 	drawHitBox         = false
	// 	inputHandler       input.InputHandler
	// 	objectAssets       assets.ObjectAssets
	// 	debugLog           debuglog.Entries
	// 	sysErrors          []error
	// 	consoleToInputChan chan console.ITxTopic
	// 	gui                ui.GUI
	// )

	// consoleToInputChan = make(chan console.ITxTopic, 1)
	// defer close(consoleToInputChan)

	// start command server
	// go console.StartServer(consoleToInputChan)
	// if Test {
	// 	go console.RunConsole()
	// }

	//setup gui
	// gui.InitGUI()

	//panic level errors
	// sysErrors = make([]error, 0)

	//load assets
	// objectAssets = rules.LoadAssets(sysErrors)

	// last := time.Now()
	// for !win.Closed() {
	// 	//handle delta
	// 	dt := time.Since(last).Seconds()
	// 	last = time.Now()

	// 	cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
	// 	win.SetMatrix(cam)

	// 	debugLog = inputHandler.HandleInput(
	// 		win,
	// 		&cam,
	// 		gameCommands,
	// 		&gameObjs,
	// 		objectAssets,
	// 		dt,
	// 		camSpeed,
	// 		&camZoom,
	// 		camZoomSpeed,
	// 		&camPos,
	// 		&drawHitBox,
	// 		consoleToInputChan,
	// 	)

	// 	var waitGroup sync.WaitGroup

	// 	//handle game updates
	// 	gui.UpdateGUI(gameCommands)
	// 	gameCommands.ExecuteCommands(&waitGroup)
	// 	waitGroup.Wait()
	// 	gameObjs.UpdateAllObjects(dt, &waitGroup)
	// 	waitGroup.Wait()

	// 	win.Clear(colornames.Black)
	// 	//draw game objects
	// 	gameObjs.DrawAllObjects(win, drawHitBox, &waitGroup, &appState)
	// 	waitGroup.Wait()

	// 	gui.DrawGUI(win, &cam)

	// 	//draw cursor based on selected object
	// 	//must be done outside of inputHandler to be the last thing drawn
	// 	if win.MouseInsideWindow() {
	// 		if !win.Pressed(pixelgl.KeyLeftControl) {
	// 			win.SetCursorVisible(false)
	// 			//setup and object to place
	// 			inputHandler.Cursor.Draw(win, pixel.IM.Moved(cam.Unproject(win.MousePosition())))
	// 		}
	// 	} else {
	// 		win.SetCursorVisible(true)
	// 	}

	// 	win.Update()

	// 	frames++
	// 	select {
	// 	case <-second.C:
	// 		win.SetTitle(fmt.Sprintf("%s | FPS: %d | GameObjects: %d", cfg.Title, frames, len(gameObjs)))
	// 		frames = 0
	// 	default:
	// 	}

	// 	for _, entry := range debugLog {
	// 		fmt.Printf("debugLog: %s", entry.GetMessage())
	// 		if entry.GetMessage() == console.Stop {
	// 			//give time for graphics stuff finish
	// 			time.Sleep(2 * time.Second)
	// 			win.Destroy()
	// 		}
	// 	}
	// }
}
