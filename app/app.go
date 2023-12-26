// This Go code is the main package for a card game application, containing the `AppRun()` function that serves as the entry point and main loop of the game. Here's a summary of the key parts of the code:
//  1. The `import` statement imports the required packages, including Pixel for 2D graphics, PixelGL for window and input handling, and various custom packages for game logic and assets.
//  2. The `AppRun()` function initializes the game by setting up the window configuration, GUI, and assets, seeding the random number generator, and starting a command server.
//  3. The function then enters a loop that runs until the window is closed. In each iteration, it performs the following tasks:
//     a. Handles delta time for smooth frame updates and animations.
//     b. Sets up the camera and applies transformations, such as zoom and pan.
//     c. Processes user input through the `inputHandler.HandleInput()` method.
//     d. Updates game objects and the GUI using a `sync.WaitGroup` to coordinate concurrency.
//     e. Clears the window and draws a border and wooden background texture.
//     f. Draws game objects, the GUI, and a custom cursor based on the selected object.
//     g. Updates the window title with the current FPS and number of game objects.
//     h. Checks for any debug log entries with a 'console.Stop' message and closes the window if found.
//
// Overall, this package manages the main game loop and coordinates the various aspects of the game, such as input handling, object updates, rendering, and window management.
package app

import (
	"fmt"
	_ "image/png"
	"math/rand"
	"sync"
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/quartermeat/card_game/app/venderController/card_game_rules"
	"github.com/quartermeat/card_game/assets"
	"github.com/quartermeat/card_game/console"
	"github.com/quartermeat/card_game/debuglog"
	"github.com/quartermeat/card_game/gamestates"
	"github.com/quartermeat/card_game/input"
	"github.com/quartermeat/card_game/objects"
	"github.com/quartermeat/card_game/observable"
	"github.com/quartermeat/card_game/ui"
)

// AppRun() is the main game function and main loop for a card game.
// It sets up the window configuration, initializes the GUI, loads assets,
// seeds the random number generator, and starts a command server.
// It then enters a loop that handles delta time, handles input, updates game objects,
// draws game objects, draws the GUI, and draws a cursor based on selected object.
// At the end of each loop it also updates the window title with FPS and number of game objects.
// Finally it checks for any debug log entries with a message of 'console.Stop' and closes the window if found.
func AppRun() {

	StateManager := gamestates.NewStateManager()
		
	//seed rng
	rand.New(rand.NewSource(time.Now().UnixNano()))

	cfg := pixelgl.WindowConfig{
		Title:       card_game_rules.APP_TITLE,
		Bounds:      card_game_rules.WINDOW_SIZE,
		VSync:       true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var (
		appState           = observable.ObservableState{}
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
		consoleToInputChan chan console.ITxTopic
		gui                ui.GUI
	)

	// Replace the path with the path to your wooden texture image
	woodenTexture, err := assets.LoadPicture(card_game_rules.BACKGROUND_IMAGE)
	if err != nil {
		panic(err)
	}

	woodenSprite := pixel.NewSprite(woodenTexture, woodenTexture.Bounds())
	consoleToInputChan = make(chan console.ITxTopic, 1)
	defer close(consoleToInputChan)

	// start command server
	go console.StartServer(consoleToInputChan)
	if Test {
		go console.RunConsole()
	}

	// setup gui
	gui.InitGUI()

	//panic level errors
	sysErrors = make([]error, 0)

	// load assets
	objectAssets = card_game_rules.LoadAssets(sysErrors)

	last := time.Now()
	
	for !win.Closed() {
		StateManager.SetCurrentState(gamestates.Init)
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

		switch(StateManager.GetCurrentState()) {
		case gamestates.Init:{
			//setup game objects
			if(gameObjs.InitGameObjects(objectAssets, &waitGroup)){
				StateManager.SetCurrentState(gamestates.Ready)
			}
		}
		default:{
			//do nothing
		}
		}		

		//handle game updates
		gui.UpdateGUI(gameCommands)
		gameCommands.ExecuteCommands(&waitGroup)
		waitGroup.Wait()
		gameObjs.UpdateAllObjects(dt, &waitGroup)
		waitGroup.Wait()
		
		win.Clear(colornames.Black)
		
		// Draw the wooden background
		scaleX := win.Bounds().W() / woodenTexture.Bounds().W()
		scaleY := win.Bounds().H() / woodenTexture.Bounds().H()
		mat := pixel.IM.Scaled(pixel.ZV, 2.5).ScaledXY(pixel.ZV, pixel.V(scaleX, scaleY))
		woodenSprite.Draw(win, mat.Moved(win.Bounds().Center()))

		//draw game objects
		gameObjs.DrawAllObjects(win, drawHitBox, &waitGroup, &appState)
		waitGroup.Wait()

		gui.DrawGUI(win, &cam)

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

		for _, entry := range debugLog {
			fmt.Printf("debugLog: %s", entry.GetMessage())
			if entry.GetMessage() == console.Stop {
				//give time for graphics stuff finish
				time.Sleep(2 * time.Second)
				win.Destroy()
			}
		}
	}
}
