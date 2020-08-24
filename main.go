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
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  false,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	go StartServer()

	//load assets
	pinkSheet, pinkAnims, pinkAnimKeys, err := loadAnimationSheet("assets/spriteSheet.png", "assets/pink_animations.csv", 32)
	coinSheet, coinFrame, err := loadCoinSheet("assets/spriteSheet.png")
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())

	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
		gameObjs     GameObjects
		frames       = 0
		second       = time.Tick(time.Second)
		drawHitBox   = false
	)

	selectedSprite := pixel.NewSprite(coinSheet, coinFrame)

	last := time.Now()
	for !win.Closed() {
		//handle delta
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.JustPressed(pixelgl.MouseButtonRight) {
			mouse := cam.Unproject(win.MousePosition())
			selectedObj, index, hit, err := gameObjs.getSelectedGameObj(mouse)
			if err != nil {
				fmt.Printf(err.Error())
			}
			if hit {
				fmt.Println("object id:", selectedObj.id, " removed")
				gameObjs = gameObjs.fastRemoveIndex(index)
			} else {
				fmt.Println("no object selected")
			}
		}

		if win.JustPressed(pixelgl.Key0) {
			selectedSprite.Set(coinSheet, coinFrame)
		}
		if win.JustPressed(pixelgl.Key1) {
			selectedSprite.Set(pinkSheet, pinkAnims[pinkAnimKeys[0]][0])
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if win.Pressed(pixelgl.KeyLeftControl) {
				mouse := cam.Unproject(win.MousePosition())
				selectedObj, _, hit, err := gameObjs.getSelectedGameObj(mouse)
				if err != nil {
					fmt.Printf(err.Error())
				}
				if hit {
					fmt.Println("object id:", selectedObj.id)
					fmt.Println("object speed:", selectedObj.attributes.speed)
					fmt.Println("object initiative:", selectedObj.attributes.initiative)
					fmt.Println("object stamina:", selectedObj.attributes.stamina)
				} else {
					fmt.Println("no object selected")
				}
			} else {
				mouse := cam.Unproject(win.MousePosition())
				gameObjs = gameObjs.addGameObject(pinkAnimKeys, pinkAnims, pinkSheet, mouse)
			}
		}

		if win.Pressed(pixelgl.MouseButtonLeft) {
			if win.Pressed(pixelgl.KeyLeftShift) {
				mouse := cam.Unproject(win.MousePosition())
				gameObjs = gameObjs.addGameObject(pinkAnimKeys, pinkAnims, pinkSheet, mouse)
			}
		}

		if win.JustPressed(pixelgl.KeyH) {
			drawHitBox = !drawHitBox
		}

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
		if win.Pressed(pixelgl.KeyLeftControl) {
			camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)
		}

		win.Clear(colornames.Black)

		//this is craziness
		var waitGroup sync.WaitGroup

		//handle updates
		gameObjs.updateAll(dt, &waitGroup)
		waitGroup.Wait()
		//handle drawing
		gameObjs.drawAll(win, drawHitBox, &waitGroup)
		waitGroup.Wait()

		if win.MouseInsideWindow() {
			win.SetCursorVisible(false)
			selectedSprite.Draw(win, pixel.IM.Moved(cam.Unproject(win.MousePosition())))

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
