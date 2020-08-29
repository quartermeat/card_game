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
	sheet       pixel.Picture
	anims       map[string][]pixel.Rect
	sprite      *pixel.Sprite
	rate        float64
	state       objectState
	counter     float64
	dir         float64
	giblet      *gibletObject
	destination pixel.Vec
	vel         pixel.Vec
	hitBox      pixel.Rect
	position    pixel.Vec
	matrix      pixel.Matrix
	attributes  livingObjAttributes
}

type livingObjAttributes struct {
	initiative float64
	speed      float64
	stamina    float64
}

func (livingObj *livingObject) Sprite() *pixel.Sprite {
	return livingObj.sprite
}

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
					switch otherObj.(type) {
					case *gibletObject:
						{
							livingObj.giblet = otherObj.(*gibletObject)
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

			if livingObj.attributes.stamina <= 0 {
				livingObj.changeState(idle)
			}
		}
	case selected:
		{
			//make idle
			livingObj.sprite.Set(livingObj.sheet, livingObj.anims["idle"][interval%len(livingObj.anims["idle"])])
			livingObj.attributes.stamina += livingObj.counter
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

	if drawHitBox || livingObj.state == selected {
		imd := imdraw.New(nil)
		imd.Color = pixel.RGB(0, 255, 0)
		imd.Push(livingObj.hitBox.Min, livingObj.hitBox.Max)
		imd.Rectangle(1)
		imd.Draw(win)
	}
	waitGroup.Done()
}
