package main

import (
	"errors"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type ObjectState int

const (
	idle ObjectState = iota
	moving
	selected_idle
	selected_moving
)

const (
	maxGameObjects = 400
	maxInitiative  = 10.0
	maxSpeed       = 100.0
	maxStamina     = 100.0
)

//NextID is the next assignable object ID
var NextID = 0

type IGameObject interface {
	ObjectName() string
	Sprite() *pixel.Sprite
	GetAssets() ObjectAssets
	getID() int
	setHitBox()
	getHitBox() pixel.Rect
	update(dt float64, gameObjects GameObjects, waitGroup *sync.WaitGroup)
	changeState(newState ObjectState)
	draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup)
	moveToPosition(position pixel.Vec)
}

//GameObjects is a slice of all the gameObjects
type GameObjects []IGameObject

func (gameObjs GameObjects) fastRemoveIndex(index int) GameObjects {
	gameObjs[index] = gameObjs[len(gameObjs)-1] // Copy last element to index i.
	gameObjs = gameObjs[:len(gameObjs)-1]       // Truncate slice.
	return gameObjs
}

func (gameObjs GameObjects) appendGameObject(newObject IGameObject) GameObjects {
	if len(gameObjs) >= maxGameObjects {
		return gameObjs
	}
	gameObjs = append(gameObjs, newObject)
	return gameObjs
}

func (gameObjs GameObjects) appendLivingObject(objectAssets ObjectAssets, position pixel.Vec) GameObjects {
	newLivingObject := createNewLivingObject(objectAssets, position)
	return gameObjs.appendGameObject(&newLivingObject)
}

func (gameObjs GameObjects) appendGibletObject(objectAssets ObjectAssets, position pixel.Vec) GameObjects {
	newGibletObject := createNewGibletObject(objectAssets, position)
	return gameObjs.appendGameObject(&newGibletObject)
}

func (gameObjs GameObjects) updateAllObjects(dt float64, waitGroup *sync.WaitGroup) {
	for _, currentObj := range gameObjs {
		waitGroup.Add(1)
		go currentObj.update(dt, gameObjs, waitGroup)
	}
}

func (gameObjs GameObjects) drawAllObjects(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	for _, obj := range gameObjs {
		waitGroup.Add(1)
		go obj.draw(win, drawHitBox, waitGroup)
	}
}

func (gameObjs GameObjects) getSelectedGameObjAtPosition(position pixel.Vec) (IGameObject, int, bool, error) {
	foundObject := true
	noIndex := -1

	if gameObjs == nil || len(gameObjs) == 0 {
		return nil, noIndex, !foundObject, errors.New("getSelectedGameObj: no game object exist")
	}
	for index, object := range gameObjs {
		if object.getHitBox().Contains(position) {
			return object, index, foundObject, nil
		}
	}
	return gameObjs[0], noIndex, !foundObject, nil
}
