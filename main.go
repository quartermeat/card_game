package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/errors"
	"golang.org/x/image/colornames"
)

func loadAnimationSheet(sheetPath, descPath string, frameWidth float64) (sheet pixel.Picture, anims map[string][]pixel.Rect, err error) {
	// total hack, nicely format the error at the end, so I don't have to type it every time
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "error loading animation sheet")
		}
	}()

	// open and load the spritesheet
	sheetFile, err := os.Open(sheetPath)
	if err != nil {
		return nil, nil, err
	}
	defer sheetFile.Close()
	sheetImg, _, err := image.Decode(sheetFile)
	if err != nil {
		return nil, nil, err
	}
	sheet = pixel.PictureDataFromImage(sheetImg)

	// create a slice of frames inside the spritesheet
	var frames []pixel.Rect
	for x := 0.0; x+frameWidth <= sheet.Bounds().Max.X; x += frameWidth {
		frames = append(frames, pixel.R(
			x,
			0,
			x+frameWidth,
			sheet.Bounds().H(),
		))
	}

	descFile, err := os.Open(descPath)
	if err != nil {
		return nil, nil, err
	}
	defer descFile.Close()

	anims = make(map[string][]pixel.Rect)

	// load the animation information, name and interval inside the spritesheet
	desc := csv.NewReader(descFile)
	for {
		anim, err := desc.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		name := anim[0]
		start, _ := strconv.Atoi(anim[1])
		end, _ := strconv.Atoi(anim[2])

		anims[name] = frames[start : end+1]
	}

	return sheet, anims, nil
}

type animState int

const (
	idle animState = iota
)

func (ga *pinkAnim) draw(t pixel.Target, phys *pinkPhys) {
	if ga.sprite == nil {
		ga.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	// draw the correct frame with the correct position and direction
	ga.sprite.Set(ga.sheet, ga.frame)
	ga.sprite.Draw(t, pixel.IM.Moved(phys.rect.Center()))
}

func (gp *pinkPhys) update(dt float64) {
	// apply controls
	switch {
	default:
		gp.vel.X = 0
	}
}

func (pa *pinkAnim) update(dt float64, phys *pinkPhys) {
	pa.counter += dt

	// determine the new animation state
	var newState animState
	switch {
	case phys.vel.Len() == 0:
		newState = idle
	}

	// reset the time counter if the state changed
	if pa.state != newState {
		pa.state = newState
		pa.counter = 0
	}

	// determine the correct animation frame
	switch pa.state {
	case idle:
		i := int(math.Floor(pa.counter / pa.rate))
		pa.frame = pa.anims["HangingOut"][i%len(pa.anims["HangingOut"])]
	}
}

type pinkPhys struct {
	runSpeed  float64
	jumpSpeed float64

	rect pixel.Rect
	vel  pixel.Vec
}

type pinkAnim struct {
	sheet pixel.Picture
	anims map[string][]pixel.Rect
	rate  float64

	state   animState
	counter float64
	dir     float64

	frame pixel.Rect

	sprite *pixel.Sprite
}

func run() {

	sheet, anims, err := loadAnimationSheet("assets/pink.png", "assets/pink_animations.csv", 12)
	if err != nil {
		panic(err)
	}

	cfg := pixelgl.WindowConfig{
		Title:  "PixelLifeGo!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	phys := &pinkPhys{
		runSpeed:  64,
		jumpSpeed: 192,
		rect:      pixel.R(0, 0, 32, 32),
	}

	anim := &pinkAnim{
		sheet: sheet,
		anims: anims,
		rate:  1.0 / 10,
		dir:   +1,
	}

	batch := pixel.NewBatch(&pixel.TrianglesData{}, sheet)

	var pinkFrames []pixel.Rect
	for x := sheet.Bounds().Min.X; x < sheet.Bounds().Max.X; x += 32 {
		for y := sheet.Bounds().Min.Y; y < sheet.Bounds().Max.Y; y += 32 {
			pinkFrames = append(pinkFrames, pixel.R(x, y, x+32, y+32))
		}
	}

	var (
		camZoom      = 1.0
		camZoomSpeed = 1.2
		camPos       = pixel.ZV
		camSpeed     = 500.0
		frames       = 0
		second       = time.Tick(time.Second)
	)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			pink := pixel.NewSprite(sheet, pinkFrames[rand.Intn(len(pinkFrames))])
			mouse := cam.Unproject(win.MousePosition())
			pink.Draw(batch, pixel.IM.Scaled(pixel.ZV, 1).Moved(mouse))
		}

		// update the physics and animation
		phys.update(dt)
		anim.update(dt, phys)

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

		//draw
		win.Clear(colornames.Black)
		batch.Draw(win)
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
