// Package objects provides an interface and implementation for a slice of game objects.
// This is the slice of all initialized objects that can be controlled and displayed on the screen.
package objects

import (
	"errors"
	"sync"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/quartermeat/card_game/assets"
	"github.com/quartermeat/card_game/observable"
)

// maxGameObjects is the maximum number of game objects that can be stored.
const maxGameObjects = 400

// NextID is the generator of a new game object ID.
var NextID = 0

// IGameObject defines the interface for an object in the game that can be displayed and controlled.
type IGameObject interface {
	ObjectName() string
	GetFSM() *StateMachine
	Sprite() *pixel.Sprite
	GetAssets() assets.IObjectAsset
	GetID() int
	SetHitBox()
	GetHitBox() pixel.Rect
	Update(dt float64, gameObjects GameObjects, waitGroup *sync.WaitGroup)
	Draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup)
	MoveToPosition(position pixel.Vec)
}

// GameObjects is a slice of all game objects.
type GameObjects []IGameObject

// FastRemoveIndex removes a game object from the GameObjects slice by its index.
func (gameObjs GameObjects) FastRemoveIndex(index int) GameObjects {
	gameObjs[index] = gameObjs[len(gameObjs)-1] // Copy last element to index i.
	gameObjs = gameObjs[:len(gameObjs)-1]       // Truncate slice.
	return gameObjs
}

// AppendGameObject appends a new object to the game objects slice.
func (gameObjs GameObjects) AppendGameObject(newObject IGameObject) GameObjects {
	if len(gameObjs) >= maxGameObjects {
		return gameObjs
	}
	gameObjs = append(gameObjs, newObject)
	return gameObjs
}

// UpdateAllObjects runs the Update method for all game objects in their own goroutine.
func (gameObjs GameObjects) UpdateAllObjects(dt float64, waitGroup *sync.WaitGroup) {
	for _, currentObj := range gameObjs {
		waitGroup.Add(1)
		go currentObj.Update(dt, gameObjs, waitGroup)
	}
}

// DrawAllObjects runs the Draw method for all game objects in their own goroutine.
func (gameObjs GameObjects) DrawAllObjects(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup, app *observable.ObservableState) {
	for _, obj := range gameObjs {
		waitGroup.Add(1)
		go obj.Draw(win, drawHitBox, waitGroup)
	}
}

// GetSelectedGameObjAtPosition checks if a mouse click intersects with a game object's hitbox.
// It returns the intersected game object, its index, and whether an object was found.
// TODO: Optimize for only objects on screen. Consider changing errors to debugLog instead.
func (gameObjs GameObjects) GetSelectedGameObjAtPosition(position pixel.Vec) (IGameObject, int, bool, error) {
	foundObject := true
	noIndex := -1

	if len(gameObjs) == 0 {
		return nil, noIndex, !foundObject, errors.New("getSelectedGameObj: no game object exists")
	}
	for index, object := range gameObjs {
		if object.GetHitBox().Contains(position) {
			return object, index, foundObject, nil
		}
	}

	return gameObjs[0], noIndex, !foundObject, nil
}
