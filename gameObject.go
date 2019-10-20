package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type animState int

const (
	idle   animState = iota
	moving           //incoming
)

const MAX_GAME_OBJECTS = 800

type gameObject struct {
	sheet    pixel.Picture
	anims    map[string][]pixel.Rect
	sprite   *pixel.Sprite
	rate     float64
	state    animState
	counter  float64
	dir      float64
	location pixel.Matrix
}

type GameObjects []*gameObject

func (gameObjs GameObjects) addGameObject(animationKeys []string, animations map[string][]pixel.Rect, sheet pixel.Picture, location pixel.Matrix) GameObjects {
	if len(gameObjs) >= MAX_GAME_OBJECTS {
		return gameObjs
	} else {
		randomAnimationKey := animationKeys[rand.Intn(len(animationKeys))]
		randomAnimationFrame := rand.Intn(len(animations[randomAnimationKey]))
		newSprite := pixel.NewSprite(sheet, animations[randomAnimationKey][randomAnimationFrame])
		newObject := &gameObject{
			sheet:    sheet,
			sprite:   newSprite,
			anims:    animations,
			state:    idle,
			rate:     1.0 / 10,
			dir:      +1,
			location: location,
		}
		return append(gameObjs, newObject)
	}
}

func (gameObj *gameObject) update(dt float64) {
	gameObj.counter += dt

	// determine the new animation state, do some check here and set
	newState := idle

	// determine the correct animation frame
	switch gameObj.state {
	case idle:
		i := int(math.Floor(gameObj.counter / gameObj.rate))
		gameObj.sprite.Set(gameObj.sheet, gameObj.anims["idle"][i%len(gameObj.anims["idle"])])
	}

	// reset the time counter if the state changed
	if gameObj.state != newState {
		gameObj.state = newState
		gameObj.counter = 0
	}
}

func (gameObj *gameObject) draw(win *pixelgl.Window) {
	gameObj.sprite.Draw(win, gameObj.location)
}

func (gameObjs GameObjects) updateAll(dt float64) {
	for _, obj := range gameObjs {
		obj.update(dt)
	}
}

func (gameObjs GameObjects) drawAll(win *pixelgl.Window) {
	for _, obj := range gameObjs {
		obj.draw(win)
	}
}
