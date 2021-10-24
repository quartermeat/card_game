package main

import (
	"fmt"
	_ "image/png"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "Aeon Ex Machina",
		Bounds: pixel.R(0, 0, 1280, 960),
		VSync:  false,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	go StartServer()

	//load assets
	pinkSheet, pinkAnims, pinkAnimKeys, err := loadCellAnimationSheet("assets/spriteSheet.png", "assets/pinkAnimations.csv", 32)
	if err != nil {
		panic(err)
	}
	gibletSheet, gibletAnims, gibletAnimKeys, err := loadgibletAnimationSheet("assets/spriteSheet.png", "assets/gibletAnimations.csv", 16)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())

	var (
		camPos         = pixel.ZV
		camSpeed       = 500.0
		camZoom        = 1.0
		camZoomSpeed   = 1.2
		gameObjs       GameObjects
		gameCommands   = make(Commands)
		selectedObject gameObject
		objectToPlace  gameObject
		frames         = 0
		second         = time.Tick(time.Second)
		drawHitBox     = false
	)

	objectToPlace = getShallowGibletObject(gibletAnimKeys, gibletAnims, gibletSheet)

	last := time.Now()
	for !win.Closed() {
		//handle delta
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		//handle removing an object
		if win.JustPressed(pixelgl.MouseButtonRight) && win.Pressed(pixelgl.KeyLeftControl) {
			mouse := cam.Unproject(win.MousePosition())
			//add a command to commands
			gameCommands[fmt.Sprintf("RemoveObject x:%f, y:%f", mouse.X, mouse.Y)] = gameObjs.RemoveObject(mouse)
		}

		//select giblet
		if win.JustPressed(pixelgl.Key0) {
			switch objectToPlace.(type) {
			case *gibletObject:
				{
					//do nothing, already selected
				}
			case *livingObject:
				{
					objectToPlace = getShallowGibletObject(gibletAnimKeys, gibletAnims, gibletSheet)
				}
			}
		}

		//select living object
		if win.JustPressed(pixelgl.Key1) {
			switch objectToPlace.(type) {
			case *gibletObject:
				{
					objectToPlace = getShallowLivingObject(pinkAnimKeys, pinkAnims, pinkSheet)
				}
			case *livingObject:
				{
					//do nothing, already selected
				}
			}
		}

		//place the selected object
		if win.JustPressed(pixelgl.MouseButtonLeft) && !win.Pressed(pixelgl.KeyLeftControl) {
			mouse := cam.Unproject(win.MousePosition())
			// once objectToPlace gets animation information, we can remove the type switch here
			gameCommands[fmt.Sprintf("AddObject: %s", objectToPlace.ObjectName())] = gameObjs.AddObject(objectToPlace, mouse)
		}

		//handle ctrl functions
		if win.Pressed(pixelgl.KeyLeftControl) {
			win.SetCursorVisible(true)
			//select an object
			if win.JustPressed(pixelgl.MouseButtonLeft) {
				mouse := cam.Unproject(win.MousePosition())
				newSelectedObject, _, hit, err := gameObjs.getSelectedGameObj(mouse)
				if err != nil {
					fmt.Print(err.Error())
				}
				if hit {
					//unselect last object
					if selectedObject != nil {
						selectedObject.changeState(idle)
					}

					selectedObject = newSelectedObject
					fmt.Println("object id:", selectedObject.getID())
					switch selectedObject.(type) {
					case *livingObject:
						{
							selectedObject.changeState(selected)
						}
					case *gibletObject:
						{

						}
					}
				} else {
					//add move command here
					fmt.Println("destination selected")
				}
			}
		}

		//toggle hit box draw
		if win.JustPressed(pixelgl.KeyH) {
			drawHitBox = !drawHitBox
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
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		//used for framerate test
		if win.Pressed(pixelgl.MouseButtonLeft) {
			if win.Pressed(pixelgl.KeyLeftShift) {
				mouse := cam.Unproject(win.MousePosition())
				switch objectToPlace.(type) {
				case *livingObject:
					{
						gameObjs = gameObjs.appendLivingObject(pinkAnimKeys, pinkAnims, pinkSheet, mouse)
					}
				case *gibletObject:
					{
						gameObjs = gameObjs.appendGibletObject(gibletAnimKeys, gibletAnims, gibletSheet, mouse)
					}
				}
			}
		}

		var waitGroup sync.WaitGroup

		//handle game updates
		gameCommands.executeCommands(&waitGroup)
		waitGroup.Wait()
		gameObjs.updateAllObjects(dt, &waitGroup)
		waitGroup.Wait()

		win.Clear(colornames.Black)
		//draw game objects
		gameObjs.drawAllObjects(win, drawHitBox, &waitGroup)
		waitGroup.Wait()

		//draw cursor based on selected object
		if win.MouseInsideWindow() {
			if !win.Pressed(pixelgl.KeyLeftControl) {
				win.SetCursorVisible(false)
				objectToPlace.Sprite().Draw(win, pixel.IM.Moved(cam.Unproject(win.MousePosition())))
			}
		} else {
			win.SetCursorVisible(true)
		}

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | CELLS: %d", cfg.Title, frames, len(gameObjs)))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
