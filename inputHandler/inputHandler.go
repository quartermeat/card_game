package input

import (
	"fmt"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/quartermeat/aeonExMachina/assets"
	objects "github.com/quartermeat/aeonExMachina/gameObjects"
)

type InputHandler struct {
	initialized   bool
	ObjectToPlace objects.IGameObject
}

func (input *InputHandler) InitializeObjectToPlace(object objects.IGameObject) {
	input.ObjectToPlace = object
	input.initialized = true
}

func (input *InputHandler) HandleInput(
	win *pixelgl.Window,
	cam *pixel.Matrix,
	gameCommands Commands,
	gameObjs *objects.GameObjects,
	gibletAssets assets.ObjectAssets,
	livingAssets assets.ObjectAssets,
	dt float64,
	camSpeed float64,
	camZoom *float64,
	camZoomSpeed float64,
	camPos *pixel.Vec,
	drawHitBox *bool,
) {
	if !input.initialized {
		input.InitializeObjectToPlace(objects.GetShallowLivingObject(livingAssets))
	}

	//select giblet
	if win.JustPressed(pixelgl.Key0) {
		input.ObjectToPlace = objects.GetShallowGibletObject(gibletAssets)
	}

	//select living object
	if win.JustPressed(pixelgl.Key1) {
		input.ObjectToPlace = objects.GetShallowLivingObject(livingAssets)
	}

	//place the selected object
	if win.JustPressed(pixelgl.MouseButtonLeft) && !win.Pressed(pixelgl.KeyLeftControl) {
		mouse := cam.Unproject(win.MousePosition())
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", mouse.X, mouse.Y, input.ObjectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, input.ObjectToPlace, mouse)
	}

	//move selected object to position
	if win.JustPressed(pixelgl.MouseButtonRight) {
		mouse := cam.Unproject(win.MousePosition())
		gameCommands[fmt.Sprintf("MoveSelectedToPosition: x:%f, y:%f", mouse.X, mouse.Y)] = MoveSelectedToPositionObject(gameObjs, mouse)
	}

	//handle ctrl functions
	if win.Pressed(pixelgl.KeyLeftControl) {
		win.SetCursorVisible(true)
		if win.JustPressed(pixelgl.MouseButtonRight) {
			mouse := cam.Unproject(win.MousePosition())
			gameCommands[fmt.Sprintf("RemoveObjectAtPosition x:%f, y:%f", mouse.X, mouse.Y)] = RemoveObjectAtPosition(gameObjs, mouse)
		}
		if win.JustPressed(pixelgl.MouseButtonLeft) { //ctrl + left click
			mouse := cam.Unproject(win.MousePosition())
			gameCommands[fmt.Sprintf("SelectObjectAtPosition x:%f, y:%f", mouse.X, mouse.Y)] = SelectObjectAtPosition(gameObjs, mouse)
		}
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

	// //used for framerate test
	if win.Pressed(pixelgl.MouseButtonLeft) {
		if win.Pressed(pixelgl.KeyLeftShift) {
			mouse := cam.Unproject(win.MousePosition())
			gameCommands[fmt.Sprintf("AddObject: %s", input.ObjectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, input.ObjectToPlace, mouse)
		}
	}

}
