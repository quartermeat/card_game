package input

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/exp/slices"

	"github.com/quartermeat/card_game/assets"
	"github.com/quartermeat/card_game/console"
	"github.com/quartermeat/card_game/errormgmt"
	"github.com/quartermeat/card_game/objects"
)

type InputHandler struct {
	initialized  bool
	Cursor       *pixel.Sprite
	CursorAssets assets.ObjectAsset
	win          *pixelgl.Window
	cam          *pixel.Matrix
	consoleInput <-chan console.IConsoleTxCommand
}

func (input *InputHandler) setCursor(pressed bool) {

	if !pressed {
		input.Cursor = pixel.NewSprite(input.CursorAssets.Sheet, input.CursorAssets.Anims[assets.CursorAnimations][0])
	} else {
		input.Cursor = pixel.NewSprite(input.CursorAssets.Sheet, input.CursorAssets.Anims[assets.CursorAnimations][1])
	}

	input.initialized = true
}

func (input *InputHandler) handleConsole(someFlag bool, errors errormgmt.Errors) errormgmt.Errors {
	select {
	case consoleCommand := <-input.consoleInput:
		{
			if consoleCommand.GetCommand() == console.Poke {
				someFlag = !someFlag
				input.setCursor(someFlag)
			}
			if consoleCommand.GetCommand() == console.Stop {
				stopCommand := errormgmt.AemError{
					Message: console.Stop,
				}
				errors = append(errors, stopCommand)
				return errors
			}
		}
	default:
		{
			//don't do anything
		}
	}
	return errors
}

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
	readConsole <-chan console.IConsoleTxCommand,
) (errors errormgmt.Errors) {
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
			indexError := errormgmt.AemError{
				Message: "{c.Description} is not in assests",
			}
			errors = append(errors, indexError)
		}
	}

	input.consoleInput = readConsole
	errors = input.handleConsole(cursorToggle, errors)

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

	//place the selected object
	if win.Pressed(pixelgl.MouseButtonLeft) && !win.Pressed(pixelgl.KeyLeftControl) {
		input.setCursor(true)
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

	return errors
}
