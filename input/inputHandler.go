package input

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/quartermeat/aeonExMachina/assets"
	"github.com/quartermeat/aeonExMachina/objects"
	"golang.org/x/exp/slices"
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
) {
	if !input.initialized {
		var idx int = 0
		idx = slices.IndexFunc(objectAssets, func(c assets.ObjectAsset) bool { return c.Description == assets.CursorAnimations })
		if idx != -1 {
			input.CursorAssets = objectAssets[idx]
			input.SetCursor(false)
		}
	}

	//select giblet
	// if win.JustPressed(pixelgl.Key0) {
	// 	// input.ObjectToPlace = objects.GetShallowGibletObject(gibletAssets)
	// }

	//select living object
	// if win.JustPressed(pixelgl.Key1) {
	// 	// input.ObjectToPlace = objects.GetShallowLivingObject(livingAssets)
	// }

	if win.JustReleased(pixelgl.MouseButtonLeft) && !win.Pressed(pixelgl.KeyLeftControl) {
		input.SetCursor(false)
	}

	//place the selected object
	if win.Pressed(pixelgl.MouseButtonLeft) && !win.Pressed(pixelgl.KeyLeftControl) {
		// mouse := cam.Unproject(win.MousePosition())
		// gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", mouse.X, mouse.Y, input.ObjectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, input.ObjectToPlace, mouse)
		input.SetCursor(true)
	}

	//move selected object to position
	// if win.JustPressed(pixelgl.MouseButtonRight) {
	// 	// mouse := cam.Unproject(win.MousePosition())
	// 	// gameCommands[fmt.Sprintf("MoveSelectedToPosition: x:%f, y:%f", mouse.X, mouse.Y)] = MoveSelectedToPositionObject(gameObjs, mouse)
	// }

	//handle ctrl functions
	// if win.Pressed(pixelgl.KeyLeftControl) {
	// 	// win.SetCursorVisible(true)
	// 	// if win.JustPressed(pixelgl.MouseButtonRight) {
	// 	// 	mouse := cam.Unproject(win.MousePosition())
	// 	// 	gameCommands[fmt.Sprintf("RemoveObjectAtPosition x:%f, y:%f", mouse.X, mouse.Y)] = RemoveObjectAtPosition(gameObjs, mouse)
	// 	// }
	// 	// if win.JustPressed(pixelgl.MouseButtonLeft) { //ctrl + left click
	// 	// 	mouse := cam.Unproject(win.MousePosition())
	// 	// 	gameCommands[fmt.Sprintf("SelectObjectAtPosition x:%f, y:%f", mouse.X, mouse.Y)] = SelectObjectAtPosition(gameObjs, mouse)
	// 	// }
	// }

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

	// //used for framerate test
	// if win.Pressed(pixelgl.MouseButtonLeft) {
	// 	if win.Pressed(pixelgl.KeyLeftShift) {
	// 		mouse := cam.Unproject(win.MousePosition())
	// 		gameCommands[fmt.Sprintf("AddObject: %s", input.ObjectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, input.ObjectToPlace, mouse)
	// 	}
	// }

}
