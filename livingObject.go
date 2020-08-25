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

const (
	maxLivingObjects = 400
)

type livingObject struct {
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
	attributes livingObjAttributes
}

type livingObjAttributes struct {
	initiative float64
	speed      float64
	stamina    float64
}

//LivingObjects is a slice of all the livingObjects
type LivingObjects []*livingObject

func creatNewLivingObject(animationKeys []string, animations map[string][]pixel.Rect, sheet pixel.Picture, position pixel.Vec) livingObject {
	randomAnimationKey := animationKeys[rand.Intn(len(animationKeys))]
	randomAnimationFrame := rand.Intn(len(animations[randomAnimationKey]))
	livingObj := livingObject{
		id:       NextID,
		sheet:    sheet,
		sprite:   pixel.NewSprite(sheet, animations[randomAnimationKey][randomAnimationFrame]),
		anims:    animations,
		rate:     1.0 / 10,
		dir:      0,
		position: position,
		vel:      pixel.V(0, 0),
		matrix:   pixel.IM.Moved(position),
		state:    idle,
		attributes: livingObjAttributes{
			initiative: 1 + rand.Float64()*(maxInitiative-1),
			speed:      1 + rand.Float64()*(maxSpeed-1),
			stamina:    1 + rand.Float64()*(maxStamina-1),
		},
	}
	livingObj.setHitBox()
	NextID++
	return livingObj
}

func (livingObj *livingObject) getID() int {
	return livingObj.id
}

func (livingObj *livingObject) setHitBox() {
	width := livingObj.sprite.Frame().Max.X - livingObj.sprite.Frame().Min.X
	height := livingObj.sprite.Frame().Max.Y - livingObj.sprite.Frame().Min.Y
	topRight := pixel.V(livingObj.position.X-(width/2), livingObj.position.Y-(height/2))
	bottomLeft := pixel.V(livingObj.position.X+(width/2), livingObj.position.Y+(width/2))
	livingObj.hitBox = pixel.R(topRight.X, topRight.Y, bottomLeft.X, bottomLeft.Y)
}

func (livingObj *livingObject) getHitBox() pixel.Rect {
	return livingObj.hitBox
}

func (livingObj *livingObject) update(dt float64, gameObjects GameObjects, waitGroup *sync.WaitGroup) {

	livingObj.counter += dt
	interval := int(math.Floor(livingObj.counter / livingObj.rate))
	// just have it start to move based on initiative
	switch livingObj.state {
	case idle:
		{
			//update idle animation
			livingObj.sprite.Set(livingObj.sheet, livingObj.anims["idle"][interval%len(livingObj.anims["idle"])])

			livingObj.attributes.stamina += livingObj.counter

			//start moving in a random direction
			if livingObj.counter >= livingObj.attributes.initiative {
				livingObj.changeState(moving)
			}
		}
	case moving:
		{
			//update moving animation
			livingObj.sprite.Set(livingObj.sheet, livingObj.anims["moving"][interval%len(livingObj.anims["moving"])])
			//invert x axis
			livingObj.vel.X = livingObj.attributes.speed * math.Sin(livingObj.dir) * -1
			livingObj.vel.Y = livingObj.attributes.speed * math.Cos(livingObj.dir)
			livingObj.matrix = livingObj.matrix.Moved(livingObj.vel.Scaled(dt))
			livingObj.position = livingObj.matrix.Project(livingObj.vel.Scaled(dt))
			livingObj.setHitBox()
			livingObj.attributes.stamina -= livingObj.counter

			//collision detection
			for _, otherObj := range gameObjects {
				if livingObj.hitBox.Intersects(otherObj.getHitBox()) && otherObj.getID() != livingObj.getID() {
					//handle collisions with other objects here
					switch otherObj.(type) {
					case *coinObject:
						{
							otherObj.changeState(moving)
						}
					default:
						{

						}
					}
				}
			}

			if livingObj.attributes.stamina <= 0 {
				livingObj.changeState(idle)
			}
		}
	}

	waitGroup.Done()
}

func (livingObj *livingObject) changeState(newState objectState) {
	livingObj.state = newState
	livingObj.counter = 0
	switch newState {
	case idle:
		{
			livingObj.matrix = pixel.IM.Moved(livingObj.position)
		}
	case moving:
		{
			livingObj.dir = float64(rand.Intn(360)) * (math.Pi / 180)
			livingObj.matrix = livingObj.matrix.Rotated(livingObj.position, livingObj.dir)
		}
	}
}

func (livingObj *livingObject) draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	livingObj.sprite.Draw(win, livingObj.matrix)

	if drawHitBox {
		imd := imdraw.New(nil)
		imd.Color = pixel.RGB(0, 255, 0)
		imd.Push(livingObj.hitBox.Min, livingObj.hitBox.Max)
		imd.Rectangle(1)
		imd.Draw(win)
	}
	waitGroup.Done()
}

//collection functions
func (livingObjs LivingObjects) fastRemoveIndexFromLivingObjects(index int) LivingObjects {
	livingObjs[index] = livingObjs[len(livingObjs)-1] // Copy last element to index i.
	livingObjs = livingObjs[:len(livingObjs)-1]       // Truncate slice.
	return livingObjs
}

func (livingObjs LivingObjects) updateAllLivingObjects(dt float64, gameObjs GameObjects, waitGroup *sync.WaitGroup) {
	for i := 0; i < len(livingObjs); i++ {
		waitGroup.Add(1)
		go livingObjs[i].update(dt, gameObjs, waitGroup)
	}
}

func (livingObjs LivingObjects) drawAllLivingObjects(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	for _, obj := range livingObjs {
		waitGroup.Add(1)
		go obj.draw(win, drawHitBox, waitGroup)
	}
}

func (livingObjs LivingObjects) appendLivingObject(gameObjs GameObjects, animationKeys []string, animations map[string][]pixel.Rect, sheet pixel.Picture, position pixel.Vec) (LivingObjects, GameObjects) {
	if len(livingObjs) >= maxLivingObjects {
		return livingObjs, gameObjs
	}
	if len(gameObjs) >= maxGameObjects {
		return livingObjs, gameObjs
	}
	newLivingObject := creatNewLivingObject(animationKeys, animations, sheet, position)
	gameObjs = gameObjs.appendGameObject(&newLivingObject)
	return append(livingObjs, &newLivingObject), gameObjs
}

func (livingObjs LivingObjects) getSelectedLivingObj(position pixel.Vec) (livingObject, int, bool, error) {
	foundObject := true
	noIndex := -1

	if livingObjs == nil {
		return livingObject{}, noIndex, !foundObject, errors.New("no game object exist")
	}
	for index, object := range livingObjs {
		if object.hitBox.Contains(position) {
			return *object, index, foundObject, nil
		}
	}
	return *livingObjs[0], noIndex, !foundObject, nil
}
