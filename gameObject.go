package main

import (
	"errors"
	"math"
	"math/rand"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type objectState int

const (
	idle   objectState = iota
	moving             //incoming
)

const (
	maxGameObjects = 400
	maxInitiative  = 10.0
	maxSpeed       = 100.0
	maxStamina     = 100.0
)

//nextID is the next assignable ID
var nextID = 0

type gameObject struct {
	id      int
	sheet   pixel.Picture
	anims   map[string][]pixel.Rect
	sprite  *pixel.Sprite
	rate    float64
	state   objectState
	counter float64
	dir     float64

	vel        pixel.Vec
	hitBox     pixel.Rect
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

func (gameObjs GameObjects) fastRemoveIndex(index int) GameObjects {
	gameObjs[index] = gameObjs[len(gameObjs)-1] // Copy last element to index i.
	gameObjs[len(gameObjs)-1] = nil             // Erase last element (write zero value).
	gameObjs = gameObjs[:len(gameObjs)-1]       // Truncate slice.
	return gameObjs
}

func (gameObjs GameObjects) getSelectedGameObj(position pixel.Vec) (gameObject, int, bool, error) {
	foundObject := true
	noIndex := -1

	if gameObjs == nil {
		return gameObject{}, noIndex, !foundObject, errors.New("no game object exist")
	}
	for index, object := range gameObjs {
		if object.hitBox.Contains(position) {
			return *object, index, foundObject, nil
		}
	}
	return *gameObjs[0], noIndex, !foundObject, nil
}

func setHitBox(position pixel.Vec, hitBox pixel.Rect) pixel.Rect {
	width := hitBox.Max.X - hitBox.Min.X
	height := hitBox.Max.Y - hitBox.Min.Y
	topRight := pixel.V(position.X-(width/2), position.Y-(height/2))
	bottomLeft := pixel.V(position.X+(width/2), position.Y+(width/2))
	return pixel.R(topRight.X, topRight.Y, bottomLeft.X, bottomLeft.Y)
}

func (gameObjs GameObjects) addGameObject(animationKeys []string, animations map[string][]pixel.Rect, sheet pixel.Picture, position pixel.Vec) GameObjects {
	if len(gameObjs) >= maxGameObjects {
		return gameObjs
	}
	randomAnimationKey := animationKeys[rand.Intn(len(animationKeys))]
	randomAnimationFrame := rand.Intn(len(animations[randomAnimationKey]))
	newSprite := pixel.NewSprite(sheet, animations[randomAnimationKey][randomAnimationFrame])
	newObject := &gameObject{
		id:       nextID,
		sheet:    sheet,
		sprite:   newSprite,
		anims:    animations,
		state:    idle,
		rate:     1.0 / 10,
		dir:      0, //direction in radians
		position: position,
		vel:      pixel.V(0, 0),
		matrix:   pixel.IM.Moved(position),
		hitBox:   setHitBox(position, newSprite.Frame()),
		attributes: objAttributes{
			initiative: 1 + rand.Float64()*(maxInitiative-1),
			speed:      1 + rand.Float64()*(maxSpeed-1),
			stamina:    1 + rand.Float64()*(maxStamina-1),
		},
	}
	nextID++
	return append(gameObjs, newObject)

}

func (gameObj *gameObject) update(dt float64, gameObjects GameObjects, waitGroup *sync.WaitGroup) {

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
				gameObj.changeState(moving)
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
			gameObj.position = gameObj.matrix.Project(gameObj.vel.Scaled(dt))
			gameObj.hitBox = setHitBox(gameObj.position, gameObj.sprite.Frame())
			gameObj.attributes.stamina -= gameObj.counter

			//collision detection
			for _, otherObj := range gameObjects {
				if gameObj.hitBox.Intersects(otherObj.hitBox) && otherObj.id != gameObj.id {
					// fmt.Println("two objects touching(", otherObj.id, ",", gameObj.id, "):", time.Now().UnixNano())
				}
			}

			if gameObj.attributes.stamina <= 0 {
				gameObj.changeState(idle)
			}
		}
	}

	waitGroup.Done()
}

func (gameObj *gameObject) changeState(newState objectState) {
	gameObj.state = newState
	gameObj.counter = 0
	switch newState {
	case idle:
		{
			gameObj.matrix = pixel.IM.Moved(gameObj.position)
		}
	case moving:
		{
			gameObj.dir = float64(rand.Intn(360)) * (math.Pi / 180)
			gameObj.matrix = gameObj.matrix.Rotated(gameObj.position, gameObj.dir)
		}
	}
}

func (gameObj *gameObject) draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	gameObj.sprite.Draw(win, gameObj.matrix)

	if drawHitBox {
		imd := imdraw.New(nil)
		imd.Color = pixel.RGB(0, 255, 0)
		imd.Push(gameObj.hitBox.Min, gameObj.hitBox.Max)
		imd.Rectangle(1)
		imd.Draw(win)
	}
	waitGroup.Done()
}

func (gameObjs GameObjects) updateAll(dt float64, waitGroup *sync.WaitGroup) {
	for _, obj := range gameObjs {
		waitGroup.Add(1)
		go obj.update(dt, gameObjs, waitGroup)
	}
}

func (gameObjs GameObjects) drawAll(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	for _, obj := range gameObjs {
		waitGroup.Add(1)
		go obj.draw(win, drawHitBox, waitGroup)
	}
}
