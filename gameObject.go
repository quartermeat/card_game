package main

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type animState int

const (
	idle   animState = iota
	moving           //incoming
)

//MaxGameObjects is the max number of gameObjects to be rendered
const MaxGameObjects = 800

type gameObject struct {
	sheet   pixel.Picture
	anims   map[string][]pixel.Rect
	sprite  *pixel.Sprite
	rate    float64
	state   animState
	counter float64
	dir     float64

	vel        pixel.Vec
	position   pixel.Vec
	matrix     pixel.Matrix
	attributes objAttributes
}

type objAttributes struct {
	initiative float64
	speed      float64
	stamina    float64
}

//GameObjects is a slice of all the gameObjects
type GameObjects []*gameObject

func (gameObjs GameObjects) getSelectedGameObj(position pixel.Vec) (gameObject, error) {
	if gameObjs == nil {
		return gameObject{}, errors.New("no game object exist")
	}
	for i := 0; i < len(gameObjs); i++ {

	}
	return *gameObjs[0], nil
}

func (gameObjs GameObjects) addGameObject(animationKeys []string, animations map[string][]pixel.Rect, sheet pixel.Picture, position pixel.Vec) GameObjects {
	if len(gameObjs) >= MaxGameObjects {
		return gameObjs
	}
	randomAnimationKey := animationKeys[rand.Intn(len(animationKeys))]
	randomAnimationFrame := rand.Intn(len(animations[randomAnimationKey]))
	newSprite := pixel.NewSprite(sheet, animations[randomAnimationKey][randomAnimationFrame])
	newObject := &gameObject{
		sheet:    sheet,
		sprite:   newSprite,
		anims:    animations,
		state:    idle,
		rate:     1.0 / 10,
		dir:      0, //direction in radians
		position: position,
		vel:      pixel.V(0, 0),
		matrix:   pixel.IM.Moved(position),
		attributes: objAttributes{
			initiative: 1 + rand.Float64()*(10-1),
			speed:      1 + rand.Float64()*(100-1),
			stamina:    1 + rand.Float64()*(100-1),
		},
	}
	return append(gameObjs, newObject)

}

var toggle bool = true

func (gameObj *gameObject) update(dt float64, waitGroup *sync.WaitGroup) {

	rand.Seed(time.Now().UnixNano())

	gameObj.counter += dt
	interval := int(math.Floor(gameObj.counter / gameObj.rate))
	// just have it start to move based on initiative
	switch gameObj.state {
	case idle:
		{
			//update idle animation
			gameObj.sprite.Set(gameObj.sheet, gameObj.anims["idle"][interval%len(gameObj.anims["idle"])])

			gameObj.attributes.stamina += gameObj.counter

			//start moving in a random direction
			if gameObj.counter >= gameObj.attributes.initiative {
				gameObj.state = moving
				gameObj.dir = float64(rand.Intn(360)) * (math.Pi / 180)
				gameObj.matrix = gameObj.matrix.Rotated(gameObj.position, gameObj.dir)
				gameObj.counter = 0
			}
		}
	case moving:
		{
			//update moving animation
			gameObj.sprite.Set(gameObj.sheet, gameObj.anims["moving"][interval%len(gameObj.anims["moving"])])
			//invert x axis
			gameObj.vel.X = gameObj.attributes.speed * math.Sin(gameObj.dir) * -1
			gameObj.vel.Y = gameObj.attributes.speed * math.Cos(gameObj.dir)
			gameObj.matrix = gameObj.matrix.Moved(gameObj.vel.Scaled(dt))

			gameObj.attributes.stamina -= gameObj.counter

			if gameObj.attributes.stamina <= 0 {
				gameObj.state = idle
				gameObj.counter = 0
				gameObj.position = gameObj.matrix.Project(gameObj.vel.Scaled(dt))
				gameObj.matrix = pixel.IM.Moved(gameObj.position)
			}
		}
	}

	waitGroup.Done()
}

func (gameObj *gameObject) draw(win *pixelgl.Window, waitGroup *sync.WaitGroup) {
	gameObj.sprite.Draw(win, gameObj.matrix)
	waitGroup.Done()
}

func (gameObjs GameObjects) updateAll(dt float64, waitGroup *sync.WaitGroup) {
	for _, obj := range gameObjs {
		waitGroup.Add(1)
		go obj.update(dt, waitGroup)
	}
}

func (gameObjs GameObjects) drawAll(win *pixelgl.Window, waitGroup *sync.WaitGroup) {
	for _, obj := range gameObjs {
		waitGroup.Add(1)
		obj.draw(win, waitGroup)
	}
}
