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
	selected
)

const (
	maxGameObjects = 400
	maxInitiative  = 10.0
	maxSpeed       = 100.0
	maxStamina     = 100.0
)

//NextID is the next assignable object ID
var NextID = 0

type gameObject interface {
	Sprite() *pixel.Sprite
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
	gameObjs = gameObjs[:len(gameObjs)-1]       // Truncate slice.
	return gameObjs
}

func (gameObjs GameObjects) appendGameObject(newObject gameObject) GameObjects {
	if len(gameObjs) >= maxGameObjects {
		return gameObjs
	}
	gameObjs = append(gameObjs, newObject)
	return gameObjs
}

func (gameObjs GameObjects) appendLivingObject(animationKeys []string, animations map[string][]pixel.Rect, sheet pixel.Picture, position pixel.Vec) GameObjects {
	newLivingObject := creatNewLivingObject(animationKeys, animations, sheet, position)
	return gameObjs.appendGameObject(&newLivingObject)
}

func (gameObjs GameObjects) appendGibletObject(animationKeys []string, animations map[string][]pixel.Rect, sheet pixel.Picture, position pixel.Vec) GameObjects {
	newGibletObject := creatNewGibletObject(animationKeys, animations, sheet, position)
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

func (gameObjs GameObjects) getSelectedGameObj(position pixel.Vec) (gameObject, int, bool, error) {
	foundObject := true
	noIndex := -1

	if gameObjs == nil || len(gameObjs) == 0 {
		return nil, noIndex, !foundObject, errors.New("no game object exist")
	}
	for index, object := range gameObjs {
		if object.getHitBox().Contains(position) {
			return object, index, foundObject, nil
		}
	}
	return gameObjs[0], noIndex, !foundObject, nil
}
