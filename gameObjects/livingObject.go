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

//LivingObject is an object with behavior
type LivingObject struct {
	id          int
	assets      assets.ObjectAssets
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
	scalingFactor := 0.75
	width := (livingObj.sprite.Frame().Max.X - livingObj.sprite.Frame().Min.X) * scalingFactor
	height := (livingObj.sprite.Frame().Max.Y - livingObj.sprite.Frame().Min.Y)
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
	// just have it start to move based on initiative
	switch livingObj.state {
	case IDLE:
		{
			//update idle animation
			angle := livingObj.dir * (180 / math.Pi)
			if angle >= 337.5 || angle < 22.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_right"][interval%len(livingObj.assets.Anims["idle_right"])])
			} else if angle >= 22.5 && angle < 67.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_45"][interval%len(livingObj.assets.Anims["idle_45"])])
			} else if angle >= 67.5 && angle < 112.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_up"][interval%len(livingObj.assets.Anims["idle_up"])])
			} else if angle >= 112.5 && angle < 157.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_135"][interval%len(livingObj.assets.Anims["idle_135"])])
			} else if angle >= 157.5 && angle < 202.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_left"][interval%len(livingObj.assets.Anims["idle_left"])])
			} else if angle >= 202.5 && angle < 247.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_225"][interval%len(livingObj.assets.Anims["idle_225"])])
			} else if angle >= 247.5 && angle < 292.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_down"][interval%len(livingObj.assets.Anims["idle_down"])])
			} else if angle >= 292.5 && angle < 337.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_315"][interval%len(livingObj.assets.Anims["idle_315"])])
			}

			livingObj.attributes.stamina += livingObj.counter

			//start moving in a random direction
			if livingObj.counter >= livingObj.attributes.initiative {
				livingObj.ChangeState(MOVING)
			}
		}
	case MOVING:
		{
			//update moving animation
			angle := livingObj.dir * (180 / math.Pi)
			if angle >= 337.5 || angle < 22.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["moving_right"][interval%len(livingObj.assets.Anims["moving_right"])])
			} else if angle >= 22.5 && angle < 67.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["moving_45"][interval%len(livingObj.assets.Anims["moving_45"])])
			} else if angle >= 67.5 && angle < 112.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["moving_up"][interval%len(livingObj.assets.Anims["moving_up"])])
			} else if angle >= 112.5 && angle < 157.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["moving_135"][interval%len(livingObj.assets.Anims["moving_135"])])
			} else if angle >= 157.5 && angle < 202.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["moving_left"][interval%len(livingObj.assets.Anims["moving_left"])])
			} else if angle >= 202.5 && angle < 247.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["moving_225"][interval%len(livingObj.assets.Anims["moving_225"])])
			} else if angle >= 247.5 && angle < 292.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["moving_down"][interval%len(livingObj.assets.Anims["moving_down"])])
			} else if angle >= 292.5 && angle < 337.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["moving_315"][interval%len(livingObj.assets.Anims["moving_315"])])
			}

			//invert x axis
			livingObj.vel.X = livingObj.attributes.speed * math.Cos(livingObj.dir)
			livingObj.vel.Y = livingObj.attributes.speed * math.Sin(livingObj.dir)
			livingObj.matrix = livingObj.matrix.Moved(livingObj.vel.Scaled(dt))
			livingObj.position = livingObj.matrix.Project(livingObj.vel.Scaled(dt))
			livingObj.SetHitBox()
			livingObj.attributes.stamina -= livingObj.counter

			//handle holding a giblet
			// if livingObj.giblet != nil {
			// 	//update giblet's position
			// 	livingObj.giblet.matrix = livingObj.matrix.Moved(livingObj.vel.Scaled(dt))
			// 	livingObj.giblet.position = livingObj.matrix.Project(livingObj.vel.Scaled(dt))
			// 	livingObj.attributes.stamina -= float64(livingObj.giblet.attributes.value)
			// }

			// //collision detection
			// for _, otherObj := range gameObjects {
			// 	if livingObj.hitBox.Intersects(otherObj.GetHitBox()) && otherObj.GetID() != livingObj.GetID() {
			// 		//handle collisions with other objects here
			// 		switch otherOject := otherObj.(type) {
			// 		case *GibletObject:
			// 			{
			// 				//FIXME: this needs to be updated
			// 				livingObj.giblet = otherOject
			// 				// if livingObj.giblet.host != nil && livingObj.giblet.host != livingObj {
			// 				// 	//take giblet from other host
			// 				// 	livingObj.giblet.host.giblet = nil
			// 				// }
			// 				livingObj.giblet.ChangeState(MOVING)
			// 				livingObj.giblet.host = livingObj
			// 			}
			// 		default:
			// 			{

			// 			}
			// 		}
			// 	}
			// }

			if livingObj.position == livingObj.destination {
				livingObj.motives.destinationReached = true
				livingObj.ChangeState(IDLE)
			}
			if livingObj.attributes.stamina <= 0 {
				livingObj.ChangeState(IDLE)
			}
		}
	case SELECTED_IDLE:
		{
			//make idle
			//update idle animation
			angle := livingObj.dir * (180 / math.Pi)
			if angle >= 337.5 || angle < 22.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_right"][interval%len(livingObj.assets.Anims["idle_right"])])
			} else if angle >= 22.5 && angle < 67.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_45"][interval%len(livingObj.assets.Anims["idle_45"])])
			} else if angle >= 67.5 && angle < 112.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_up"][interval%len(livingObj.assets.Anims["idle_up"])])
			} else if angle >= 112.5 && angle < 157.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_135"][interval%len(livingObj.assets.Anims["idle_135"])])
			} else if angle >= 157.5 && angle < 202.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_left"][interval%len(livingObj.assets.Anims["idle_left"])])
			} else if angle >= 202.5 && angle < 247.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_225"][interval%len(livingObj.assets.Anims["idle_225"])])
			} else if angle >= 247.5 && angle < 292.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_down"][interval%len(livingObj.assets.Anims["idle_down"])])
			} else if angle >= 292.5 && angle < 337.5 {
				livingObj.sprite.Set(livingObj.assets.Sheet, livingObj.assets.Anims["idle_315"][interval%len(livingObj.assets.Anims["idle_315"])])
			}
			livingObj.attributes.stamina += livingObj.counter
		}
	}

	waitGroup.Done()
}

func (livingObj *LivingObject) ChangeState(newState ObjectState) {
	// if livingObj.state != newState {
	// 	livingObj.lastState = livingObj.state
	// 	livingObj.state = newState
	// }
	livingObj.state = newState
	livingObj.counter = 0
	switch newState {
	case IDLE:
		{
			livingObj.matrix = pixel.IM.Moved(livingObj.position)
			angle := livingObj.dir * (180 / math.Pi)
			if angle >= 337.5 || angle < 22.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_right"][0])
			} else if angle >= 22.5 && angle < 67.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_45"][0])
			} else if angle >= 67.5 && angle < 112.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_up"][0])
			} else if angle >= 112.5 && angle < 157.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_135"][0])
			} else if angle >= 157.5 && angle < 202.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_left"][0])
			} else if angle >= 202.5 && angle < 247.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_225"][0])
			} else if angle >= 247.5 && angle < 292.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_down"][0])
			} else if angle >= 292.5 && angle < 337.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_315"][0])
			}
		}
	case MOVING:
		{
			if livingObj.motives.destinationReached {
				randFloat := float64(rand.Intn(360))
				livingObj.dir = randFloat * (math.Pi / 180)
			}
			angle := livingObj.dir * (180 / math.Pi)
			if angle >= 337.5 || angle < 22.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["moving_right"][0])
			} else if angle >= 22.5 && angle < 67.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["moving_45"][0])
			} else if angle >= 67.5 && angle < 112.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["moving_up"][0])
			} else if angle >= 112.5 && angle < 157.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["moving_135"][0])
			} else if angle >= 157.5 && angle < 202.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["moving_left"][0])
			} else if angle >= 202.5 && angle < 247.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["moving_225"][0])
			} else if angle >= 247.5 && angle < 292.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["moving_down"][0])
			} else if angle >= 292.5 && angle < 337.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["moving_315"][0])
			}

		}
	case SELECTED_IDLE:
		{
			livingObj.matrix = pixel.IM.Moved(livingObj.position)
			angle := livingObj.dir * (180 / math.Pi)
			if angle >= 337.5 || angle < 22.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_right"][0])
			} else if angle >= 22.5 && angle < 67.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_45"][0])
			} else if angle >= 67.5 && angle < 112.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_up"][0])
			} else if angle >= 112.5 && angle < 157.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_135"][0])
			} else if angle >= 157.5 && angle < 202.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_left"][0])
			} else if angle >= 202.5 && angle < 247.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_225"][0])
			} else if angle >= 247.5 && angle < 292.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_down"][0])
			} else if angle >= 292.5 && angle < 337.5 {
				livingObj.sprite = pixel.NewSprite(livingObj.assets.Sheet, livingObj.assets.Anims["idle_315"][0])
			}
		}
	case SELECTED_MOVING:
		{
			livingObj.motives.destinationReached = false
			livingObj.ChangeState(MOVING)
		}
	}
}

func (livingObj *LivingObject) GetState() ObjectState {
	return livingObj.state
}

func (livingObj *LivingObject) Draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	livingObj.sprite.Draw(win, livingObj.matrix)

	if drawHitBox || livingObj.state == SELECTED_IDLE {
		imd := imdraw.New(nil)
		imd.Color = pixel.RGB(0, 255, 0)
		imd.Push(livingObj.hitBox.Min, livingObj.hitBox.Max)
		imd.Rectangle(1)
		imd.Draw(win)
	}
	waitGroup.Done()
}

func (livingObj *LivingObject) MoveToPosition(position pixel.Vec) {
	livingObj.destination = position
	livingObj.dir = livingObj.position.To(livingObj.destination).Angle()
	if livingObj.dir < 0 {
		livingObj.dir = livingObj.dir + (2 * math.Pi)
	}
	livingObj.ChangeState(SELECTED_MOVING)
}

//#endregion

func GetShallowLivingObject(objectAssets assets.ObjectAssets) *LivingObject {
	return &LivingObject{
		id:       -1,
		assets:   objectAssets,
		sprite:   pixel.NewSprite(objectAssets.Sheet, objectAssets.Anims["idle_down"][0]),
		rate:     1.0 / 2,
		dir:      270 * (math.Pi / 180),
		position: pixel.V(0, 0),
		vel:      pixel.V(0, 0),
		giblet:   nil,
		matrix:   pixel.IM.Moved(pixel.ZV),
		state:    IDLE,
		attributes: livingObjAttributes{
			initiative: 0,
			speed:      0,
			stamina:    0,
		},
	}
}

func createNewLivingObject(objectAssets assets.ObjectAssets, position pixel.Vec) LivingObject {
	randomAnimationKey := objectAssets.AnimKeys[rand.Intn(len(objectAssets.AnimKeys))]
	randomAnimationFrame := rand.Intn(len(objectAssets.Anims[randomAnimationKey]))
	livingObj := LivingObject{
		id:       NextID,
		assets:   objectAssets,
		sprite:   pixel.NewSprite(objectAssets.Sheet, objectAssets.Anims[randomAnimationKey][randomAnimationFrame]),
		rate:     1.0 / 10,
		dir:      270 * (math.Pi / 180),
		giblet:   nil,
		position: position,
		vel:      pixel.V(0, 0),
		matrix:   pixel.IM.Moved(position),
		state:    IDLE,
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
