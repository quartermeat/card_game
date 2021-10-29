package objects

import (
	"math"
	"math/rand"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/quartermeat/aeonExMachina/assets"
)

type ILivingState interface {
	EnterState()
	GoMoving(obj *LivingObject)
	DoMoving(obj *LivingObject)
	GoIdle(obj *LivingObject)
	DoIdle(obj *LivingObject, interval int)
	GetState() string
}

//LivingObject is an object with behavior
type LivingObject struct {
	id      int
	assets  assets.ObjectAssets
	sprite  *pixel.Sprite
	rate    float64
	state   ILivingState
	counter float64
	dir     float64
	// giblet      *GibletObject
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

func (livingObj *LivingObject) ObjectName() string {
	return "Living"
}

func (livingObj *LivingObject) Sprite() *pixel.Sprite {
	return livingObj.sprite
}

func (livingObj *LivingObject) GetAssets() assets.ObjectAssets {
	return livingObj.assets
}

func (livingObj *LivingObject) GetID() int {
	return livingObj.id
}

func (livingObj *LivingObject) SetHitBox() {
	width := livingObj.sprite.Frame().Max.X - livingObj.sprite.Frame().Min.X
	height := livingObj.sprite.Frame().Max.Y - livingObj.sprite.Frame().Min.Y
	topRight := pixel.V(livingObj.position.X-(width/2), livingObj.position.Y-(height/2))
	bottomLeft := pixel.V(livingObj.position.X+(width/2), livingObj.position.Y+(width/2))
	livingObj.hitBox = pixel.R(topRight.X, topRight.Y, bottomLeft.X, bottomLeft.Y)
}

func (livingObj *LivingObject) GetHitBox() pixel.Rect {
	return livingObj.hitBox
}

func (livingObj *LivingObject) Update(dt float64, gameObjects GameObjects, waitGroup *sync.WaitGroup) {

	livingObj.counter += dt
	interval := int(math.Floor(livingObj.counter / livingObj.rate))
	livingObj.state.DoIdle(livingObj, interval)
	// // just have it start to move based on initiative
	// switch livingObj.state {
	// case IDLE:
	// 	{
	// 		//update idle animation
	// 		livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle"][interval%len(livingObj.assets.Anims["idle"])])

	// 		livingObj.attributes.stamina += livingObj.counter

	// 		//start moving in a random direction
	// 		if livingObj.counter >= livingObj.attributes.initiative {
	// 			livingObj.ChangeState(MOVING)
	// 		}
	// 	}
	// case MOVING:
	// 	{
	// 		//update moving animation
	// 		livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["moving"][interval%len(livingObj.assets.Anims["moving"])])
	// 		//invert x axis
	// 		livingObj.vel.X = livingObj.attributes.speed * math.Cos(livingObj.dir)
	// 		livingObj.vel.Y = livingObj.attributes.speed * math.Sin(livingObj.dir)
	// 		livingObj.matrix = livingObj.matrix.Moved(livingObj.vel.Scaled(dt))
	// 		livingObj.position = livingObj.matrix.Project(livingObj.vel.Scaled(dt))
	// 		livingObj.SetHitBox()
	// 		livingObj.attributes.stamina -= livingObj.counter

	// 		//handle holding a giblet
	// 		if livingObj.giblet != nil {
	// 			//update giblet's position
	// 			livingObj.giblet.matrix = livingObj.matrix.Moved(livingObj.vel.Scaled(dt))
	// 			livingObj.giblet.position = livingObj.matrix.Project(livingObj.vel.Scaled(dt))
	// 			livingObj.attributes.stamina -= float64(livingObj.giblet.attributes.value)
	// 		}

	// 		//collision detection
	// 		for _, otherObj := range gameObjects {
	// 			if livingObj.hitBox.Intersects(otherObj.GetHitBox()) && otherObj.GetID() != livingObj.GetID() {
	// 				//handle collisions with other objects here
	// 				switch otherOject := otherObj.(type) {
	// 				case *GibletObject:
	// 					{
	// 						//FIXME: this needs to be updated
	// 						livingObj.giblet = otherOject
	// 						// if livingObj.giblet.host != nil && livingObj.giblet.host != livingObj {
	// 						// 	//take giblet from other host
	// 						// 	livingObj.giblet.host.giblet = nil
	// 						// }
	// 						livingObj.giblet.ChangeState(MOVING)
	// 						livingObj.giblet.host = livingObj
	// 					}
	// 				default:
	// 					{

	// 					}
	// 				}
	// 			}
	// 		}

	// 		if livingObj.position == livingObj.destination {
	// 			livingObj.motives.destinationReached = true
	// 			livingObj.ChangeState(IDLE)
	// 		}
	// 		if livingObj.attributes.stamina <= 0 {
	// 			livingObj.ChangeState(IDLE)
	// 		}
	// 	}
	// case SELECTED_IDLE:
	// 	{
	// 		//make idle
	// 		livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle"][interval%len(livingObj.assets.Anims["idle"])])
	// 		livingObj.attributes.stamina += livingObj.counter
	// 	}
	// case SELECTED_MOVING:
	// 	{
	// 		//change state to moving
	// 	}
	// }

	waitGroup.Done()
}

func (livingObj *LivingObject) ChangeControlState(newState IControlState) {

}

func (livingObj *LivingObject) ChangeState(newState ILivingState) {
	livingObj.state = newState
	livingObj.state.EnterState()
}

func (livingObj *LivingObject) Draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	livingObj.sprite.Draw(win, livingObj.matrix)

	if drawHitBox || livingObj.state.GetState() == "idle" {
		imd := imdraw.New(nil)
		imd.Color = pixel.RGB(0, 255, 0)
		imd.Push(livingObj.hitBox.Min, livingObj.hitBox.Max)
		imd.Rectangle(1)
		imd.Draw(win)
	}
	waitGroup.Done()
}

func (livingObj *LivingObject) MoveToPosition(position pixel.Vec) {
	// livingObj.destination = position
	// livingObj.dir = livingObj.position.To(livingObj.destination).Angle()
	// livingObj.ChangeState(SELECTED_MOVING)
}

//#endregion

func GetShallowLivingObject(objectAssets assets.ObjectAssets) *LivingObject {
	return &LivingObject{
		id:       -1,
		assets:   objectAssets,
		sprite:   pixel.NewSprite(objectAssets.Sheet, objectAssets.Anims["idle"][0]),
		rate:     1.0 / 2,
		dir:      0.0,
		position: pixel.V(0, 0),
		vel:      pixel.V(0, 0),
		// giblet:   nil,
		matrix: pixel.IM.Moved(pixel.ZV),
		state:  new(LivingIdleState),
		attributes: livingObjAttributes{
			initiative: 0,
			speed:      0,
			stamina:    0,
		},
	}
}

func CreateNewLivingObject(objectAssets assets.ObjectAssets, position pixel.Vec) LivingObject {
	randomAnimationKey := objectAssets.AnimKeys[rand.Intn(len(objectAssets.AnimKeys))]
	randomAnimationFrame := rand.Intn(len(objectAssets.Anims[randomAnimationKey]))
	livingObj := LivingObject{
		id:     NextID,
		assets: objectAssets,
		sprite: pixel.NewSprite(objectAssets.Sheet, objectAssets.Anims[randomAnimationKey][randomAnimationFrame]),
		rate:   1.0 / 10,
		dir:    0.0,
		// giblet:   nil,
		position: position,
		vel:      pixel.V(0, 0),
		matrix:   pixel.IM.Moved(position),
		state:    new(LivingIdleState),
		attributes: livingObjAttributes{
			initiative: 1 + rand.Float64()*(maxInitiative-1),
			speed:      1 + rand.Float64()*(maxSpeed-1),
			stamina:    1 + rand.Float64()*(maxStamina-1),
		},
		motives: livingObjMotives{
			destinationReached: true,
		},
	}
	livingObj.SetHitBox()
	NextID++
	return livingObj
}
