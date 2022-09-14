// Package 'object' provides the interface for a slice of game objects.
// This is the slice of all initialized objects that can be controlled and displayed on screen
package objects

import (
	"errors"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/quartermeat/card_game/assets"
)

// MAX GAME OBJECTS
const (
	maxGameObjects = 400
)

// NextID is the generator of a new game object ID
var NextID = 0

// IGameObject holds the capabilities of an object in the game that can be displayed and controlled
type IGameObject interface {
	ObjectName() string
	GetFSM() *StateMachine
	Sprite() *pixel.Sprite
	GetAssets() assets.ObjectAsset
	GetID() int
	SetHitBox()
	GetHitBox() pixel.Rect
	Update(dt float64, gameObjects GameObjects, waitGroup *sync.WaitGroup)
	Draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup)
	MoveToPosition(position pixel.Vec)
}

// GameObjects is a slice of all the gameObjects
type GameObjects []IGameObject

// FastRemoveIndex removes a gameObject from the GameObjects slice by it's index
func (gameObjs GameObjects) FastRemoveIndex(index int) GameObjects {
	gameObjs[index] = gameObjs[len(gameObjs)-1] // Copy last element to index i.
	gameObjs = gameObjs[:len(gameObjs)-1]       // Truncate slice.
	return gameObjs
}

// AppendGameObject appends a new object to the game objects slice
func (gameObjs GameObjects) AppendGameObject(newObject IGameObject) GameObjects {
	if len(gameObjs) >= maxGameObjects {
		return gameObjs
	}
	gameObjs = append(gameObjs, newObject)
	return gameObjs
}

// UpdateAllObjects runs all game objects Update method within it's own go routine
func (gameObjs GameObjects) UpdateAllObjects(dt float64, waitGroup *sync.WaitGroup) {
	for _, currentObj := range gameObjs {
		waitGroup.Add(1)
		go currentObj.Update(dt, gameObjs, waitGroup)
	}
}

// DrawAllObjects runs all game objects Draw method within it's own go routine
func (gameObjs GameObjects) DrawAllObjects(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	for _, obj := range gameObjs {
		waitGroup.Add(1)
		go obj.Draw(win, drawHitBox, waitGroup)
	}
}

// GetSelectedGameObjAtPosition intersects a mouse click with any game objects hitbox
// TODO: maybe optimize for only objects on screen
// TODO: change to debugLog instead of error
func (gameObjs GameObjects) GetSelectedGameObjAtPosition(position pixel.Vec) (IGameObject, int, bool, error) {
	foundObject := true
	noIndex := -1

	if len(gameObjs) == 0 {
		return nil, noIndex, !foundObject, errors.New("getSelectedGameObj: no game object exist")
	}
	for index, object := range gameObjs {
		if object.GetHitBox().Contains(position) {
			return object, index, foundObject, nil
		}
	}

	return gameObjs[0], noIndex, !foundObject, nil
}

// business defined logic///////////////////////////////////////////////////
// Get Card returns the card from the game object slice based on the index
// func (gameObjs GameObjects) GetCard(idx int) (card.ICard, error) {

// }
