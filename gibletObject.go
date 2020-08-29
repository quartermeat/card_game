package main

import (
	"math"
	"math/rand"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type gibletObject struct {
	id         int
	sheet      pixel.Picture
	anims      map[string][]pixel.Rect
	sprite     *pixel.Sprite
	rate       float64
	state      objectState
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

func (gibletObj *gibletObject) Sprite() *pixel.Sprite {
	return gibletObj.sprite
}

func (gibletObj *gibletObject) getID() int {
	return gibletObj.id
}

func (gibletObj *gibletObject) setHitBox() {
	width := gibletObj.sprite.Frame().Max.X - gibletObj.sprite.Frame().Min.X
	height := gibletObj.sprite.Frame().Max.Y - gibletObj.sprite.Frame().Min.Y
	topRight := pixel.V(gibletObj.position.X-(width/2), gibletObj.position.Y-(height/2))
	bottomLeft := pixel.V(gibletObj.position.X+(width/2), gibletObj.position.Y+(width/2))
	gibletObj.hitBox = pixel.R(topRight.X, topRight.Y, bottomLeft.X, bottomLeft.Y)
}

func (gibletObj *gibletObject) getHitBox() pixel.Rect {
	return gibletObj.hitBox
}

func creatNewGibletObject(animationKeys []string, animations map[string][]pixel.Rect, sheet pixel.Picture, position pixel.Vec) gibletObject {
	randomAnimationKey := animationKeys[rand.Intn(len(animationKeys))]
	randomAnimationFrame := rand.Intn(len(animations[randomAnimationKey]))
	gibletObj := gibletObject{
		id:       NextID,
		sheet:    sheet,
		sprite:   pixel.NewSprite(sheet, animations[randomAnimationKey][randomAnimationFrame]),
		anims:    animations,
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

func (gibletObj *gibletObject) changeState(newState objectState) {
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

func (gibletObj *gibletObject) update(dt float64, gameObjects GameObjects, waitGroup *sync.WaitGroup) {
	gibletObj.counter += dt
	interval := int(math.Floor(gibletObj.counter / gibletObj.rate))
	// just have it start to move based on initiative
	switch gibletObj.state {
	case idle: //lying on ground
		{
			//update idle animation
			gibletObj.sprite.Set(gibletObj.sheet, gibletObj.anims["gibletIdle"][interval%len(gibletObj.anims["gibletIdle"])])
		}
	case moving: //held
		{
			//update moving animation
			gibletObj.sprite.Set(gibletObj.sheet, gibletObj.anims["gibletIdle"][interval%len(gibletObj.anims["gibletIdle"])])
		}
	}

	waitGroup.Done()
}

func (gibletObj *gibletObject) draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
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
