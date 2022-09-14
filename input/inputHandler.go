package input

import (
	"fmt"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/exp/slices"

	"github.com/quartermeat/card_game/assets"
	"github.com/quartermeat/card_game/console"
	"github.com/quartermeat/card_game/debuglog"
	"github.com/quartermeat/card_game/domainObjects/card"
	"github.com/quartermeat/card_game/objects"
)

// InputHandler is a monolithic struct to handle user interactions with the app
type InputHandler struct {
	initialized  bool
	Cursor       *pixel.Sprite
	CursorAssets assets.ObjectAsset
	win          *pixelgl.Window
	cam          *pixel.Matrix
	consoleInput <-chan console.ITxTopic
}

func (input *InputHandler) setCursor(pressed bool) {

	if !pressed {
		input.Cursor = pixel.NewSprite(input.CursorAssets.Sheet, input.CursorAssets.Anims[assets.CursorAnimations][0])
	} else {
		input.Cursor = pixel.NewSprite(input.CursorAssets.Sheet, input.CursorAssets.Anims[assets.CursorAnimations][1])
	}

	input.initialized = true
}

func (input *InputHandler) handleConsole(someFlag bool, debugLog debuglog.Entries) debuglog.Entries {
	select {
	case consoleCommand := <-input.consoleInput:
		{
			if consoleCommand.GetTopicId() == console.Poke {
				someFlag = !someFlag
				input.setCursor(someFlag)
			}
			if consoleCommand.GetTopicId() == console.Stop {
				stopCommand := debuglog.Entry{
					Message: console.Stop,
				}
				debugLog = append(debugLog, stopCommand)
				return debugLog
			}
		}
	default:
		{
			//don't do anything
		}
	}
	return debugLog
}

// HandleInput is a super method ran from main
// atm: handles input from the keyboard, mouse and console
func (input *InputHandler) HandleInput(
	win *pixelgl.Window,
	cam *pixel.Matrix,
	gameCommands Commands,
	gameObjs *objects.GameObjects,
	objectAssets assets.ObjectAssets,
	dt float64,
	camSpeed float64,
	camZoom *float64,
	camZoomSpeed float64,
	camPos *pixel.Vec,
	drawHitBox *bool,
	readConsole <-chan console.ITxTopic,
) (debugLog debuglog.Entries) {
	//defaults
	var (
		cursorToggle bool
	)

	//do initialization of input handler
	if !input.initialized {
		//set window and cam
		input.win = win
		input.cam = cam

		//set cursor
		cursorToggle = false
		var idx int = 0
		idx = slices.IndexFunc(objectAssets, func(c assets.ObjectAsset) bool { return c.Description == assets.CursorAnimations })
		if idx != -1 {
			input.CursorAssets = objectAssets[idx]
			input.setCursor(cursorToggle)
		} else {
			indexError := debuglog.Entry{
				Message: "{c.Description} is not in assests",
			}
			debugLog = append(debugLog, indexError)
		}
	}

	input.consoleInput = readConsole
	debugLog = input.handleConsole(cursorToggle, debugLog)

	if win.MouseInsideWindow() {
		if !win.Pressed(pixelgl.KeyLeftControl) {
			win.SetCursorVisible(false)
			//setup and object to place
			input.Cursor.Draw(win, pixel.IM.Moved(cam.Unproject(win.MousePosition())))
		}
	} else {
		win.SetCursorVisible(true)
	}

	if win.JustReleased(pixelgl.MouseButtonLeft) && !win.Pressed(pixelgl.KeyLeftControl) {
		input.setCursor(false)
	}

	//handle ctrl functions
	if win.Pressed(pixelgl.KeyLeftControl) {
		input.setCursor(true)
		win.SetCursorVisible(true)
		if win.JustPressed(pixelgl.MouseButtonLeft) { //ctrl + left click
			mouse := cam.Unproject(win.MousePosition())
			gameCommands[fmt.Sprintf("SelectObjectAtPosition x:%f, y:%f", mouse.X, mouse.Y)] = SelectObjectAtPosition(gameObjs, mouse)
		}
	}

	//place the selected object
	if win.Pressed(pixelgl.MouseButtonLeft) && !win.Pressed(pixelgl.KeyLeftControl) {
		input.setCursor(true)
	}

	if win.JustPressed(pixelgl.Key0) {
		mouse := cam.Unproject(win.MousePosition())
		objectToPlace := card.NewCardObject(objectAssets[1], mouse)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", mouse.X, mouse.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, mouse)
	}

	//toggle global hit box draw for debugging
	if win.JustPressed(pixelgl.KeyH) {
		*drawHitBox = !*drawHitBox
	}

	//move camera
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

	//zoom camera
	*camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

	return debugLog
}
