package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type animState int

const (
	idle animState = iota
	moving
)

type pinkObject struct {
	sheet   pixel.Picture
	anims   map[string][]pixel.Rect
	sprite  *pixel.Sprite
	rate    float64
	state   animState
	counter float64
	dir     float64
}

func (po *pinkObject) update(dt float64) {
	po.counter += dt

	// determine the new animation state
	var newState animState
	newState = idle

	// determine the correct animation frame
	switch po.state {
	case idle:
		i := int(math.Floor(po.counter / po.rate))
		po.sprite.Set(po.sheet, po.anims["idle"][i%len(po.anims["idle"])])
	}

	// reset the time counter if the state changed
	if po.state != newState {
		po.state = newState
		po.counter = 0
	}
}

func run() {

	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
		pinks        []*pinkObject
		matrices     []pixel.Matrix
		frames       = 0
		second       = time.Tick(time.Second)
	)

	pinkSheet, pinkAnims, pinkAnimKeys, err := loadAnimationSheet("assets/pink.png", "assets/pink_animations.csv", 32)

	if err != nil {
		panic(err)
	}

	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Life in Go",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	last := time.Now()
	for !win.Closed() {
		//handle delta
		dt := time.Since(last).Seconds()
		last = time.Now()

		//handle input
		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			randomAnimationKey := pinkAnimKeys[rand.Intn(len(pinkAnimKeys))]
			randomAnimationFrame := rand.Intn(len(pinkAnims[randomAnimationKey]))
			pinkSprite := pixel.NewSprite(pinkSheet, pinkAnims[randomAnimationKey][randomAnimationFrame])
			pinkObject := &pinkObject{
				sheet:  pinkSheet,
				sprite: pinkSprite,
				anims:  pinkAnims,
				state:  idle,
				rate:   1.0 / 10,
				dir:    +1,
			}
			pinks = append(pinks, pinkObject)
			mouse := cam.Unproject(win.MousePosition())
			matrices = append(matrices, pixel.IM.Scaled(pixel.ZV, 1).Moved(mouse))
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
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		//handle updates
		for _, pink := range pinks {
			pink.update(dt)
		}

		//handle drawing
		win.Clear(colornames.Black)

		for i, pink := range pinks {
			pink.sprite.Draw(win, matrices[i])
		}

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | CELLS: %d", cfg.Title, frames, len(pinks)))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
