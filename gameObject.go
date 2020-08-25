package main

import (
	"errors"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type objectState int

const (
	idle   objectState = iota
	moving             //incoming
)

const (
	maxGameObjects = 800
	maxInitiative  = 10.0
	maxSpeed       = 100.0
	maxStamina     = 100.0
)

//NextID is the next assignable object ID
var NextID = 0

type gameObject interface {
	getID() int
	setHitBox()
	getHitBox() pixel.Rect
	update(dt float64, gameObjects GameObjects, waitGroup *sync.WaitGroup)
	changeState(newState objectState)
	draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup)
}

//GameObjects is a slice of all the gameObjects
type GameObjects []gameObject

func (gameObjs GameObjects) fastRemoveIndex(index int) GameObjects {
	gameObjs[index] = gameObjs[len(gameObjs)-1] // Copy last element to index i.
	gameObjs[len(gameObjs)-1] = nil             // Erase last element (write zero value).
	gameObjs = gameObjs[:len(gameObjs)-1]       // Truncate slice.
	return gameObjs
}

func (gameObjs GameObjects) appendGameObject(newObject gameObject) GameObjects {
	if len(gameObjs) >= maxGameObjects {
		return nil
	}
	gameObjs = append(gameObjs, newObject)
	return gameObjs
}

func (gameObjs GameObjects) getSelectedGameObj(position pixel.Vec) (gameObject, int, bool, error) {
	foundObject := true
	noIndex := -1

	if gameObjs == nil {
		return nil, noIndex, !foundObject, errors.New("no game object exist")
	}
	for index, object := range gameObjs {
		if object.getHitBox().Contains(position) {
			return object, index, foundObject, nil
		}
	}
	return gameObjs[0], noIndex, !foundObject, nil
}
