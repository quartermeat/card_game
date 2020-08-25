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
	maxCoinObjects = 100
)

type coinObject struct {
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
	attributes coinObjAttributes
}

type coinObjAttributes struct {
	value int
}

//CoinObjects is a slice of all the livingObjects
type CoinObjects []*coinObject

func (coinObj coinObject) getID() int {
	return coinObj.id
}

func (coinObj *coinObject) setHitBox() {
	width := coinObj.sprite.Frame().Max.X - coinObj.sprite.Frame().Min.X
	height := coinObj.sprite.Frame().Max.Y - coinObj.sprite.Frame().Min.Y
	topRight := pixel.V(coinObj.position.X-(width/2), coinObj.position.Y-(height/2))
	bottomLeft := pixel.V(coinObj.position.X+(width/2), coinObj.position.Y+(width/2))
	coinObj.hitBox = pixel.R(topRight.X, topRight.Y, bottomLeft.X, bottomLeft.Y)
}

func (coinObj coinObject) getHitBox() pixel.Rect {
	return coinObj.hitBox
}

func creatNewCoinObject(animationKeys []string, animations map[string][]pixel.Rect, sheet pixel.Picture, position pixel.Vec) coinObject {
	randomAnimationKey := animationKeys[rand.Intn(len(animationKeys))]
	randomAnimationFrame := rand.Intn(len(animations[randomAnimationKey]))
	coinObj := coinObject{
		id:       NextID,
		sheet:    sheet,
		sprite:   pixel.NewSprite(sheet, animations[randomAnimationKey][randomAnimationFrame]),
		anims:    animations,
		rate:     1.0 / 2,
		dir:      0,
		position: position,
		vel:      pixel.V(0, 0),
		matrix:   pixel.IM.Moved(position),
		state:    idle,
		attributes: coinObjAttributes{
			value: 1,
		},
	}
	coinObj.setHitBox()
	NextID++
	return coinObj
}

func (coinObj *coinObject) changeState(newState objectState) {
	coinObj.state = newState
	coinObj.counter = 0
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

func (coinObj *coinObject) update(dt float64, gameObjects GameObjects, waitGroup *sync.WaitGroup) {
	coinObj.counter += dt
	interval := int(math.Floor(coinObj.counter / coinObj.rate))
	// just have it start to move based on initiative
	switch coinObj.state {
	case idle: //lying on ground
		{
			//update idle animation
			coinObj.sprite.Set(coinObj.sheet, coinObj.anims["coinIdle"][interval%len(coinObj.anims["coinIdle"])])
		}
	case moving: //held
		{
			//update moving animation
			coinObj.sprite.Set(coinObj.sheet, coinObj.anims["coinIdle"][interval%len(coinObj.anims["coinIdle"])])
			//add logic to move with host
		}
	}

	waitGroup.Done()
}

func (coinObj coinObject) draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	coinObj.sprite.Draw(win, coinObj.matrix)

	if drawHitBox {
		imd := imdraw.New(nil)
		imd.Color = pixel.RGB(0, 255, 0)
		imd.Push(coinObj.hitBox.Min, coinObj.hitBox.Max)
		imd.Rectangle(1)
		imd.Draw(win)
	}
	waitGroup.Done()
}

//collection functions
func (coinObjs CoinObjects) fastRemoveIndexFromCoinObjects(index int) CoinObjects {
	coinObjs[index] = coinObjs[len(coinObjs)-1] // Copy last element to index i.
	coinObjs = coinObjs[:len(coinObjs)-1]       // Truncate slice.
	return coinObjs
}

func (coinObjs CoinObjects) updateAllCoinObjects(dt float64, gameObjs GameObjects, waitGroup *sync.WaitGroup) {
	for i := 0; i < len(coinObjs); i++ {
		waitGroup.Add(1)
		go coinObjs[i].update(dt, gameObjs, waitGroup)
	}
}

func (coinObjs CoinObjects) drawAllCoinObjects(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	for _, obj := range coinObjs {
		waitGroup.Add(1)
		go obj.draw(win, drawHitBox, waitGroup)
	}
}

func (coinObjs CoinObjects) appendCoinObject(gameObjs GameObjects, animationKeys []string, animations map[string][]pixel.Rect, sheet pixel.Picture, position pixel.Vec) (CoinObjects, GameObjects) {
	if len(coinObjs) >= maxCoinObjects {
		return coinObjs, gameObjs
	}
	if len(gameObjs) >= maxGameObjects {
		return coinObjs, gameObjs
	}
	newCoinObject := creatNewCoinObject(animationKeys, animations, sheet, position)
	gameObjs = gameObjs.appendGameObject(&newCoinObject)
	return append(coinObjs, &newCoinObject), gameObjs
}

func (coinObjs CoinObjects) getSelectedCoinObj(position pixel.Vec) (coinObject, int, bool, error) {
	foundObject := true
	noIndex := -1

	if coinObjs == nil {
		return coinObject{}, noIndex, !foundObject, errors.New("no game object exist")
	}
	for index, object := range coinObjs {
		if object.hitBox.Contains(position) {
			return *object, index, foundObject, nil
		}
	}
	return *coinObjs[0], noIndex, !foundObject, nil
}
