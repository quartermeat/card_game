package objects

import (
	"errors"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/quartermeat/aeonExMachina/assets"
)

const (
	maxGameObjects = 400
	maxInitiative  = 10.0
	maxSpeed       = 100.0
	maxStamina     = 100.0
)

type IControlState interface {
	EnterState()
	Select(IGameObject)
	Unselect(IGameObject)
}

//NextID is the next assignable object ID
var NextID = 0

//IGameObject is what every object should be able to do
type IGameObject interface {
	ObjectName() string
	Sprite() *pixel.Sprite
	GetAssets() assets.ObjectAssets
	GetID() int
	SetHitBox()
	GetHitBox() pixel.Rect
	Update(dt float64, gameObjects GameObjects, waitGroup *sync.WaitGroup)
	ChangeControlState(newState IControlState)
	Draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup)
	MoveToPosition(position pixel.Vec)
}

//GameObjects is a slice of all the gameObjects
type GameObjects []IGameObject

func (gameObjs GameObjects) FastRemoveIndex(index int) GameObjects {
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

func (gameObjs GameObjects) AppendLivingObject(objectAssets assets.ObjectAssets, position pixel.Vec) GameObjects {
	newLivingObject := CreateNewLivingObject(objectAssets, position)
	return gameObjs.appendGameObject(&newLivingObject)
}

// func (gameObjs GameObjects) AppendGibletObject(objectAssets assets.ObjectAssets, position pixel.Vec) GameObjects {
// 	newGibletObject := createNewGibletObject(objectAssets, position)
// 	return gameObjs.appendGameObject(&newGibletObject)
// }

func (gameObjs GameObjects) UpdateAllObjects(dt float64, waitGroup *sync.WaitGroup) {
	for _, currentObj := range gameObjs {
		waitGroup.Add(1)
		go currentObj.Update(dt, gameObjs, waitGroup)
	}
}

func (gameObjs GameObjects) DrawAllObjects(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	for _, obj := range gameObjs {
		waitGroup.Add(1)
		go obj.Draw(win, drawHitBox, waitGroup)
	}
}

func (gameObjs GameObjects) GetSelectedGameObjAtPosition(position pixel.Vec) (IGameObject, int, bool, error) {
	foundObject := true
	noIndex := -1

	if gameObjs == nil || len(gameObjs) == 0 {
		return nil, noIndex, !foundObject, errors.New("getSelectedGameObj: no game object exist")
	}
	for index, object := range gameObjs {
		if object.GetHitBox().Contains(position) {
			return object, index, foundObject, nil
		}
	}
	return gameObjs[0], noIndex, !foundObject, nil
}
