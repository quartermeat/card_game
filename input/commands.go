// Package 'input' is the interface available to make changes to game objects
package input

import (
	"fmt"
	"sync"

	"github.com/faiface/pixel"
	"github.com/quartermeat/card_game/assets"
	"github.com/quartermeat/card_game/objects"
	"github.com/quartermeat/card_game/objects/domainObjects/card"
)

// Commands is the map of commands to execute
type Commands map[string]ICommand

// ICommand interface is used to execute game commands
type ICommand interface {
	execute(*sync.WaitGroup)
}

// ExecuteCommands executes the queued list of commands
func (commands Commands) ExecuteCommands(waitGroup *sync.WaitGroup) {
	for key, value := range commands {
		fmt.Printf("executing: %s\n", key)
		waitGroup.Add(1)
		go value.execute(waitGroup)
		delete(commands, key)
	}
}

type addObjectAtPositionCommand struct {
	gameObjs      *objects.GameObjects
	objectToPlace objects.IGameObject
	position      pixel.Vec
	objectAssets  assets.IObjectAsset
}

func (command *addObjectAtPositionCommand) execute(waitGroup *sync.WaitGroup) {
	switch command.objectToPlace.(type) {
	case card.ICard:
		{
			*command.gameObjs = command.gameObjs.AppendGameObject(command.objectToPlace)
		}
	}

	waitGroup.Done()
}

// AddObjectAtPosition allows for the addition of a game object
func AddObjectAtPosition(objs *objects.GameObjects, newObject objects.IGameObject, newPosition pixel.Vec) ICommand {
	return &addObjectAtPositionCommand{
		gameObjs:      objs,
		position:      newPosition,
		objectToPlace: newObject,
		objectAssets:  newObject.GetAssets(),
	}
}

type removeObjectAtPositionCommand struct {
	gameObjs *objects.GameObjects
	position pixel.Vec
}

func (command *removeObjectAtPositionCommand) execute(waitGroup *sync.WaitGroup) {
	_, index, hit, err := command.gameObjs.GetSelectedGameObjAtPosition(command.position)
	if err != nil {
		panic(err)
	}
	if hit {
		*command.gameObjs = command.gameObjs.FastRemoveIndex(index)
	}

	waitGroup.Done()
}

// RemoveObjectAtPosition allows for the removal of a game Object based on Vec location
func RemoveObjectAtPosition(objs *objects.GameObjects, fromPosition pixel.Vec) ICommand {
	return &removeObjectAtPositionCommand{
		gameObjs: objs,
		position: fromPosition,
	}
}

type selectObjectAtPositionCommand struct {
	gameObjs *objects.GameObjects
	position pixel.Vec
}

func (command *selectObjectAtPositionCommand) execute(waitGroup *sync.WaitGroup) {
	selectedObject, _, _, err := command.gameObjs.GetSelectedGameObjAtPosition(command.position)
	if err != nil {
		panic(err)
	}

	switch selectedObject.(type) {
	case card.ICard:
		{
			err := selectedObject.GetFSM().SendEvent(card.FlipUp, selectedObject)
			if err != nil {
				_ = selectedObject.GetFSM().SendEvent(card.FlipDown, selectedObject)
			}
		}
	}

	waitGroup.Done()
}

// SelectObjectAtPosition allows for changing the state of a game object based on Vec location to selected
func SelectObjectAtPosition(objs *objects.GameObjects, fromPosition pixel.Vec) ICommand {
	return &selectObjectAtPositionCommand{
		gameObjs: objs,
		position: fromPosition,
	}
}

type moveSelectedObjectToPositionCommand struct {
	gameObjs *objects.GameObjects
	position pixel.Vec
}

func (command *moveSelectedObjectToPositionCommand) execute(waitGroup *sync.WaitGroup) {
	// for _, obj := range *command.gameObjs {
	// 	if obj.GetState() == objects.SELECTED_IDLE {
	// 		obj.MoveToPosition(command.position)
	// 	}
	// }
	waitGroup.Done()
}

// MoveSelectedObject allows for directing selected objects to move to a position
func MoveSelectedToPositionObject(objs *objects.GameObjects, fromPosition pixel.Vec) ICommand {
	return &moveSelectedObjectToPositionCommand{
		gameObjs: objs,
		position: fromPosition,
	}
}
