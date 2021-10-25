package main

import (
	"fmt"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type InputHandler struct {
	initialized   bool
	objectToPlace IGameObject
}

func (input *InputHandler) InitializeObjectToPlace(object IGameObject) {
	input.objectToPlace = object
	input.initialized = true
}

func (input *InputHandler) HandleInput(
	win *pixelgl.Window,
	cam *pixel.Matrix,
	gameCommands Commands,
	gameObjs *GameObjects,
	gibletAssets ObjectAssets,
	livingAssets ObjectAssets,
	dt float64,
	camSpeed float64,
	camZoom *float64,
	camZoomSpeed float64,
	camPos *pixel.Vec,
	drawHitBox *bool,
) {
	if !input.initialized {
		input.InitializeObjectToPlace(getShallowLivingObject(livingAssets))
	}

	//select giblet
	if win.JustPressed(pixelgl.Key0) {
		input.objectToPlace = getShallowGibletObject(gibletAssets)
	}

	//select living object
	if win.JustPressed(pixelgl.Key1) {
		input.objectToPlace = getShallowLivingObject(livingAssets)
	}

	//place the selected object
	if win.JustPressed(pixelgl.MouseButtonLeft) && !win.Pressed(pixelgl.KeyLeftControl) {
		mouse := cam.Unproject(win.MousePosition())
		// once objectToPlace gets animation information, we can remove the type switch here
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", mouse.X, mouse.Y, input.objectToPlace.ObjectName())] = gameObjs.AddObjectAtPosition(input.objectToPlace, mouse)
	}

	//move selected object to position
	if win.JustPressed(pixelgl.MouseButtonRight) {
		mouse := cam.Unproject(win.MousePosition())
		gameCommands[fmt.Sprintf("MoveSelectedToPosition: x:%f, y:%f", mouse.X, mouse.Y)] = gameObjs.MoveSelectedToPositionObject(mouse)
	}

	//handle ctrl functions
	if win.Pressed(pixelgl.KeyLeftControl) {
		win.SetCursorVisible(true)
		if win.JustPressed(pixelgl.MouseButtonRight) {
			mouse := cam.Unproject(win.MousePosition())
			//add a command to commands
			gameCommands[fmt.Sprintf("RemoveObjectAtPosition x:%f, y:%f", mouse.X, mouse.Y)] = gameObjs.RemoveObjectAtPosition(mouse)
		}
		if win.JustPressed(pixelgl.MouseButtonLeft) { //ctrl + left click
			mouse := cam.Unproject(win.MousePosition())
			gameCommands[fmt.Sprintf("SelectObjectAtPosition x:%f, y:%f", mouse.X, mouse.Y)] = gameObjs.SelectObjectAtPosition(mouse)
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
			gameCommands[fmt.Sprintf("AddObject: %s", input.objectToPlace.ObjectName())] = gameObjs.AddObjectAtPosition(input.objectToPlace, mouse)
		}
	}

}
