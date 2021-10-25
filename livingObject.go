package main

import (
	"math"
	"math/rand"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type livingObject struct {
	id          int
	assets      ObjectAssets
	sprite      *pixel.Sprite
	rate        float64
	state       ObjectState
	counter     float64
	dir         float64
	giblet      *GibletObject
	destination pixel.Vec
	vel         pixel.Vec
	hitBox      pixel.Rect
	position    pixel.Vec
	matrix      pixel.Matrix
	attributes  livingObjAttributes
	motives     livingObjMotives
}

type livingObjAttributes struct {
	initiative float64
	speed      float64
	stamina    float64
}

type livingObjMotives struct {
	destinationReached bool
}

//#region GAMEOBJECT IMPLEMENTATION

func (livingObj *livingObject) ObjectName() string {
	return "Living"
}

func (livingObj *livingObject) Sprite() *pixel.Sprite {
	return livingObj.sprite
}

func (livingObj *livingObject) GetAssets() ObjectAssets {
	return livingObj.assets
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
			livingObj.sprite.Set(livingObj.assets.sheet, livingObj.assets.anims["idle"][interval%len(livingObj.assets.anims["idle"])])

			livingObj.attributes.stamina += livingObj.counter

			//start moving in a random direction
			if livingObj.counter >= livingObj.attributes.initiative {
				livingObj.changeState(moving)
			}
		}
	case moving:
		{
			//update moving animation
			livingObj.sprite.Set(livingObj.assets.sheet, livingObj.assets.anims["moving"][interval%len(livingObj.assets.anims["moving"])])
			//invert x axis
			livingObj.vel.X = livingObj.attributes.speed * math.Sin(livingObj.dir) * -1
			livingObj.vel.Y = livingObj.attributes.speed * math.Cos(livingObj.dir)
			livingObj.matrix = livingObj.matrix.Moved(livingObj.vel.Scaled(dt))
			//ADD check destination reached is true
			//check if destination has been reached, turn off destination flag
			livingObj.position = livingObj.matrix.Project(livingObj.vel.Scaled(dt))
			livingObj.setHitBox()
			livingObj.attributes.stamina -= livingObj.counter

			//handle holding a giblet
			if livingObj.giblet != nil {
				//update giblet's position
				livingObj.giblet.matrix = livingObj.matrix.Moved(livingObj.vel.Scaled(dt))
				livingObj.giblet.position = livingObj.matrix.Project(livingObj.vel.Scaled(dt))
				livingObj.attributes.stamina -= float64(livingObj.giblet.attributes.value)
			}

			//collision detection
			for _, otherObj := range gameObjects {
				if livingObj.hitBox.Intersects(otherObj.getHitBox()) && otherObj.getID() != livingObj.getID() {
					//handle collisions with other objects here
					switch otherOject := otherObj.(type) {
					case *GibletObject:
						{
							//FIXME: this needs to be updated
							livingObj.giblet = otherOject
							// if livingObj.giblet.host != nil && livingObj.giblet.host != livingObj {
							// 	//take giblet from other host
							// 	livingObj.giblet.host.giblet = nil
							// }
							livingObj.giblet.changeState(moving)
							livingObj.giblet.host = livingObj
						}
					default:
						{

						}
					}
				}
			}

			//ADD:
			//if destination reached is true
			//changeState to idle

			if livingObj.attributes.stamina <= 0 {
				livingObj.changeState(idle)
			}
		}
	case selected_idle:
		{
			//make idle
			livingObj.sprite.Set(livingObj.assets.sheet, livingObj.assets.anims["idle"][interval%len(livingObj.assets.anims["idle"])])
			livingObj.attributes.stamina += livingObj.counter
		}
	case selected_moving:
		{
			//change state to moving
		}
	}

	waitGroup.Done()
}

func (livingObj *livingObject) changeState(newState ObjectState) {
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
	case selected_idle:
		{
			livingObj.matrix = pixel.IM.Moved(livingObj.position)
		}
	case selected_moving:
		{
			livingObj.matrix = livingObj.matrix.Rotated(livingObj.position, livingObj.dir)
			livingObj.motives.destinationReached = false
		}
	}
}

func (livingObj *livingObject) draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	livingObj.sprite.Draw(win, livingObj.matrix)

	if drawHitBox || livingObj.state == selected_idle {
		imd := imdraw.New(nil)
		imd.Color = pixel.RGB(0, 255, 0)
		imd.Push(livingObj.hitBox.Min, livingObj.hitBox.Max)
		imd.Rectangle(1)
		imd.Draw(win)
	}
	waitGroup.Done()
}

func (livingObj *livingObject) moveToPosition(position pixel.Vec) {
	livingObj.destination = position
	//calculate angle based on current location vs destination
	//update object direction based on angle

	livingObj.changeState(selected_moving)
}

//#endregion

func getShallowLivingObject(objectAssets ObjectAssets) *livingObject {
	return &livingObject{
		id:       -1,
		assets:   objectAssets,
		sprite:   pixel.NewSprite(objectAssets.sheet, objectAssets.anims["idle"][0]),
		rate:     1.0 / 2,
		dir:      0,
		position: pixel.V(0, 0),
		vel:      pixel.V(0, 0),
		giblet:   nil,
		matrix:   pixel.IM.Moved(pixel.V(0, 0)),
		state:    idle,
		attributes: livingObjAttributes{
			initiative: 0,
			speed:      0,
			stamina:    0,
		},
	}
}

func createNewLivingObject(objectAssets ObjectAssets, position pixel.Vec) livingObject {
	randomAnimationKey := objectAssets.animKeys[rand.Intn(len(objectAssets.animKeys))]
	randomAnimationFrame := rand.Intn(len(objectAssets.anims[randomAnimationKey]))
	livingObj := livingObject{
		id:       NextID,
		assets:   objectAssets,
		sprite:   pixel.NewSprite(objectAssets.sheet, objectAssets.anims[randomAnimationKey][randomAnimationFrame]),
		rate:     1.0 / 10,
		dir:      0,
		giblet:   nil,
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
