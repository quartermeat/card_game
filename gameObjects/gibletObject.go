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

type GibletObject struct {
	id         int
	assets     assets.ObjectAssets
	sprite     *pixel.Sprite
	rate       float64
	state      ObjectState
	counter    float64
	dir        float64
	host       *LivingObject
	vel        pixel.Vec
	hitBox     pixel.Rect
	position   pixel.Vec
	matrix     pixel.Matrix
	attributes gibletObjAttributes
}

type gibletObjAttributes struct {
	value int
}

//#region GAME OBJECT IMPLEMENTATION

func (gibletObj *GibletObject) ObjectName() string {
	return "Giblet"
}

func (gibletObj *GibletObject) Sprite() *pixel.Sprite {
	return gibletObj.sprite
}

func (gibletObj *GibletObject) GetAssets() assets.ObjectAssets {
	return gibletObj.assets
}

func (gibletObj *GibletObject) GetID() int {
	return gibletObj.id
}

func (gibletObj *GibletObject) SetHitBox() {
	width := gibletObj.sprite.Frame().Max.X - gibletObj.sprite.Frame().Min.X
	height := gibletObj.sprite.Frame().Max.Y - gibletObj.sprite.Frame().Min.Y
	topRight := pixel.V(gibletObj.position.X-(width/2), gibletObj.position.Y-(height/2))
	bottomLeft := pixel.V(gibletObj.position.X+(width/2), gibletObj.position.Y+(width/2))
	gibletObj.hitBox = pixel.R(topRight.X, topRight.Y, bottomLeft.X, bottomLeft.Y)
}

func (gibletObj *GibletObject) GetHitBox() pixel.Rect {
	return gibletObj.hitBox
}

func (gibletObj *GibletObject) Update(dt float64, gameObjects GameObjects, waitGroup *sync.WaitGroup) {
	gibletObj.counter += dt
	interval := int(math.Floor(gibletObj.counter / gibletObj.rate))
	// just have it start to move based on initiative
	switch gibletObj.state {
	case IDLE: //lying on ground
		{
			//update idle animation
			gibletObj.sprite.Set(gibletObj.assets.Sheet, gibletObj.assets.Anims["idle_right"][interval%len(gibletObj.assets.Anims["idle_right"])])
		}
	case MOVING: //held
		{
			//update moving animation
			gibletObj.sprite.Set(gibletObj.assets.Sheet, gibletObj.assets.Anims["moving_right"][interval%len(gibletObj.assets.Anims["moving_right"])])
		}
	}

	waitGroup.Done()
}

func (gibletObj *GibletObject) ChangeState(newState ObjectState) {
	gibletObj.state = newState
	gibletObj.counter = 0
	switch newState {
	case IDLE:
		{
			//do transistion to idle stuff
		}
	case MOVING:
		{
			//do transistion to moving stuff

		}
	}
}

func (gibletObj *GibletObject) GetState() ObjectState {
	return gibletObj.state
}

func (gibletObj *GibletObject) Draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	gibletObj.sprite.Draw(win, gibletObj.matrix)

	if drawHitBox {
		imd := imdraw.New(nil)
		imd.Color = pixel.RGB(0, 255, 0)
		imd.Push(gibletObj.hitBox.Min, gibletObj.hitBox.Max)
		imd.Rectangle(1)
		imd.Draw(win)
	}
	waitGroup.Done()
}

func (gibletObj *GibletObject) MoveToPosition(position pixel.Vec) {
	//doesn't move
}

//#endregion

func GetShallowGibletObject(objectAssets assets.ObjectAssets) *GibletObject {
	return &GibletObject{
		id:       -1,
		assets:   objectAssets,
		sprite:   pixel.NewSprite(objectAssets.Sheet, objectAssets.Anims["idle_right"][0]),
		rate:     1.0 / 2,
		dir:      0,
		position: pixel.V(0, 0),
		host:     nil,
		vel:      pixel.V(0, 0),
		matrix:   pixel.IM.Moved(pixel.V(0, 0)),
		state:    IDLE,
		attributes: gibletObjAttributes{
			value: 1,
		},
	}
}

func createNewGibletObject(objectAssets assets.ObjectAssets, position pixel.Vec) GibletObject {
	randomAnimationKey := objectAssets.AnimKeys[rand.Intn(len(objectAssets.AnimKeys))]
	randomAnimationFrame := rand.Intn(len(objectAssets.Anims[randomAnimationKey]))
	gibletObj := GibletObject{
		id:       NextID,
		assets:   objectAssets,
		sprite:   pixel.NewSprite(objectAssets.Sheet, objectAssets.Anims[randomAnimationKey][randomAnimationFrame]),
		rate:     1.0 / 2,
		dir:      0,
		position: position,
		host:     nil,
		vel:      pixel.V(0, 0),
		matrix:   pixel.IM.Moved(position),
		state:    IDLE,
		attributes: gibletObjAttributes{
			value: 1,
		},
	}
	gibletObj.SetHitBox()
	NextID++
	return gibletObj
}
