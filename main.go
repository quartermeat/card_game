package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	_ "image/png"

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
		livingObjs     LivingObjects
		gibletObjs     GibletObjects
		selectedObject gameObject
		objectToPlace  gameObject
		frames         = 0
		second         = time.Tick(time.Second)
		drawHitBox     = false
	)

	objectToPlace = &gibletObject{sprite: pixel.NewSprite(gibletSheet, gibletAnims["gibletIdle"][0])}

	last := time.Now()
	for !win.Closed() {
		//handle delta
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		//handle removing an object
		//removeObject(win, &cam, gameObjs, livingObjs, gibletObjs)
		if win.JustPressed(pixelgl.MouseButtonRight) {
			mouse := cam.Unproject(win.MousePosition())
			selectedObj, index, hit, err := gameObjs.getSelectedGameObj(mouse)
			if err != nil {
				fmt.Printf(err.Error())
			}
			if hit {
				fmt.Println("object id:", selectedObj.getID(), " removed")
				gameObjs = gameObjs.fastRemoveIndex(index)

				switch selectedObj.(type) {
				case *livingObject:
					{
						livingObjs = livingObjs.fastRemoveIndexFromLivingObjects(index)
						selectedObj = nil
					}
				case *gibletObject:
					{
						gibletObjs = gibletObjs.fastRemoveIndexFromGibletObjects(index)
						selectedObj = nil
					}
				}
			} else {
				fmt.Println("no object selected")
			}
			hit = false
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
					objectToPlace = &gibletObject{sprite: pixel.NewSprite(gibletSheet, gibletAnims[gibletAnimKeys[0]][0])}
				}
			}
		}
		//select living object
		if win.JustPressed(pixelgl.Key1) {
			switch objectToPlace.(type) {
			case *gibletObject:
				{
					objectToPlace = &livingObject{sprite: pixel.NewSprite(pinkSheet, pinkAnims[pinkAnimKeys[0]][0])}
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
			//add object based on selectedObj
			switch objectToPlace.(type) {
			case *livingObject:
				{
					livingObjs, gameObjs = livingObjs.appendLivingObject(gameObjs, pinkAnimKeys, pinkAnims, pinkSheet, mouse)
				}
			case *gibletObject:
				{
					gibletObjs, gameObjs = gibletObjs.appendGibletObject(gameObjs, gibletAnimKeys, gibletAnims, gibletSheet, mouse)
				}
			}
		}

		//handle ctrl functions
		if win.Pressed(pixelgl.KeyLeftControl) {
			win.SetCursorVisible(true)
			//select an object
			if win.JustPressed(pixelgl.MouseButtonLeft) {
				mouse := cam.Unproject(win.MousePosition())
				newSelectedObject, _, hit, err := gameObjs.getSelectedGameObj(mouse)
				if err != nil {
					fmt.Printf(err.Error())
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

		////used for framerate test
		// if win.Pressed(pixelgl.MouseButtonLeft) {
		// 	if win.Pressed(pixelgl.KeyLeftShift) {
		// 		mouse := cam.Unproject(win.MousePosition())
		// 		switch selectedObject {
		// 		case living:
		// 			{
		// 				livingObjs, gameObjs = livingObjs.appendLivingObject(gameObjs, pinkAnimKeys, pinkAnims, pinkSheet, mouse)
		// 			}
		// 		case giblet:
		// 			{
		// 				gibletObjs, gameObjs = gibletObjs.appendGibletObject(gameObjs, gibletAnimKeys, gibletAnims, gibletSheet, mouse)
		// 			}
		// 		}
		// 	}
		// }

		//this is craziness
		var waitGroup sync.WaitGroup

		win.Clear(colornames.Black)

		//handle drawing
		livingObjs.updateAllLivingObjects(dt, &gameObjs, &waitGroup)
		waitGroup.Wait()
		gibletObjs.updateAllgibletObjects(dt, &gameObjs, &waitGroup)
		waitGroup.Wait()

		livingObjs.drawAllLivingObjects(win, drawHitBox, &waitGroup)
		waitGroup.Wait()
		gibletObjs.drawAllGibletObjects(win, drawHitBox, &waitGroup)
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
