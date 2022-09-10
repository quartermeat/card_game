package input

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/exp/slices"

	"github.com/quartermeat/aeonExMachina/assets"
	"github.com/quartermeat/aeonExMachina/errormgmt"
	"github.com/quartermeat/aeonExMachina/objects"
)

type InputHandler struct {
	initialized  bool
	Cursor       *pixel.Sprite
	CursorAssets assets.ObjectAsset
}

func (input *InputHandler) SetCursor(pressed bool) {

	if !pressed {
		input.Cursor = pixel.NewSprite(input.CursorAssets.Sheet, input.CursorAssets.Anims[assets.CursorAnimations][0])
	} else {
		input.Cursor = pixel.NewSprite(input.CursorAssets.Sheet, input.CursorAssets.Anims[assets.CursorAnimations][1])
	}

	input.initialized = true
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
) (errors errormgmt.Errors) {
	//do initialization of input handler
	if !input.initialized {

		//set cursor
		var idx int = 0
		idx = slices.IndexFunc(objectAssets, func(c assets.ObjectAsset) bool { return c.Description == assets.CursorAnimations })
		if idx != -1 {
			input.CursorAssets = objectAssets[idx]
			input.SetCursor(false)
		} else {
			indexError := errormgmt.AemError{
				Message: "{c.Description} is not in assests",
			}
			errors = append(errors, indexError)
		}

	}

	if win.JustReleased(pixelgl.MouseButtonLeft) && !win.Pressed(pixelgl.KeyLeftControl) {
		input.SetCursor(false)
	}

	//place the selected object
	if win.Pressed(pixelgl.MouseButtonLeft) && !win.Pressed(pixelgl.KeyLeftControl) {
		input.SetCursor(true)
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
