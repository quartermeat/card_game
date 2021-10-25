package main

import (
	"math"
	"math/rand"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type GibletObject struct {
	id         int
	assets     ObjectAssets
	sprite     *pixel.Sprite
	rate       float64
	state      ObjectState
	counter    float64
	dir        float64
	host       *livingObject
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

func (gibletObj *GibletObject) GetAssets() ObjectAssets {
	return gibletObj.assets
}

func (gibletObj *GibletObject) getID() int {
	return gibletObj.id
}

func (gibletObj *GibletObject) setHitBox() {
	width := gibletObj.sprite.Frame().Max.X - gibletObj.sprite.Frame().Min.X
	height := gibletObj.sprite.Frame().Max.Y - gibletObj.sprite.Frame().Min.Y
	topRight := pixel.V(gibletObj.position.X-(width/2), gibletObj.position.Y-(height/2))
	bottomLeft := pixel.V(gibletObj.position.X+(width/2), gibletObj.position.Y+(width/2))
	gibletObj.hitBox = pixel.R(topRight.X, topRight.Y, bottomLeft.X, bottomLeft.Y)
}

func (gibletObj *GibletObject) getHitBox() pixel.Rect {
	return gibletObj.hitBox
}

func (gibletObj *GibletObject) update(dt float64, gameObjects GameObjects, waitGroup *sync.WaitGroup) {
	gibletObj.counter += dt
	interval := int(math.Floor(gibletObj.counter / gibletObj.rate))
	// just have it start to move based on initiative
	switch gibletObj.state {
	case idle: //lying on ground
		{
			//update idle animation
			gibletObj.sprite.Set(gibletObj.assets.sheet, gibletObj.assets.anims["gibletIdle"][interval%len(gibletObj.assets.anims["gibletIdle"])])
		}
	case moving: //held
		{
			//update moving animation
			gibletObj.sprite.Set(gibletObj.assets.sheet, gibletObj.assets.anims["gibletIdle"][interval%len(gibletObj.assets.anims["gibletIdle"])])
		}
	}

	waitGroup.Done()
}

func (gibletObj *GibletObject) changeState(newState ObjectState) {
	gibletObj.state = newState
	gibletObj.counter = 0
	switch newState {
	case idle:
		{
			//do transistion to idle stuff
		}
	case moving:
		{
			//do transistion to moving stuff

		}
	}
}

func (gibletObj *GibletObject) draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
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

func (gibletObj *GibletObject) moveToPosition(position pixel.Vec) {
	//doesn't move
}

//#endregion

func getShallowGibletObject(objectAssets ObjectAssets) *GibletObject {
	return &GibletObject{
		id:       -1,
		assets:   objectAssets,
		sprite:   pixel.NewSprite(objectAssets.sheet, objectAssets.anims["gibletIdle"][0]),
		rate:     1.0 / 2,
		dir:      0,
		position: pixel.V(0, 0),
		host:     nil,
		vel:      pixel.V(0, 0),
		matrix:   pixel.IM.Moved(pixel.V(0, 0)),
		state:    idle,
		attributes: gibletObjAttributes{
			value: 1,
		},
	}
}

func createNewGibletObject(objectAssets ObjectAssets, position pixel.Vec) GibletObject {
	randomAnimationKey := objectAssets.animKeys[rand.Intn(len(objectAssets.animKeys))]
	randomAnimationFrame := rand.Intn(len(objectAssets.anims[randomAnimationKey]))
	gibletObj := GibletObject{
		id:       NextID,
		assets:   objectAssets,
		sprite:   pixel.NewSprite(objectAssets.sheet, objectAssets.anims[randomAnimationKey][randomAnimationFrame]),
		rate:     1.0 / 2,
		dir:      0,
		position: position,
		host:     nil,
		vel:      pixel.V(0, 0),
		matrix:   pixel.IM.Moved(position),
		state:    idle,
		attributes: gibletObjAttributes{
			value: 1,
		},
	}
	gibletObj.setHitBox()
	NextID++
	return gibletObj
}
