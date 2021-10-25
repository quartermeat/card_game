package main

import (
	"fmt"
	"sync"

	"github.com/faiface/pixel"
)

//Commands is the map of commands to execute
type Commands map[string]ICommand

//ICommand interface is used to execute game commands
type ICommand interface {
	execute(*sync.WaitGroup)
}

//concurrently execute queued game commands
func (commands Commands) executeCommands(waitGroup *sync.WaitGroup) {
	for key, value := range commands {
		fmt.Printf("executing: %s\n", key)
		waitGroup.Add(1)
		go value.execute(waitGroup)
		delete(commands, key)
	}
}

//#region ADD OBJECT COMMAND

type addObjectAtPositionCommand struct {
	gameObjs      *GameObjects
	objectToPlace IGameObject
	position      pixel.Vec
	objectAssets  ObjectAssets
}

func (command *addObjectAtPositionCommand) execute(waitGroup *sync.WaitGroup) {
	switch command.objectToPlace.(type) {
	case *livingObject:
		{
			*command.gameObjs = command.gameObjs.appendLivingObject(command.objectAssets, command.position)
		}
	case *GibletObject:
		{
			*command.gameObjs = command.gameObjs.appendGibletObject(command.objectAssets, command.position)
		}
	}

	waitGroup.Done()
}

//AddObjectAtPosition allows for the addition of a game object
func (objects *GameObjects) AddObjectAtPosition(newObject IGameObject, newPosition pixel.Vec) ICommand {
	return &addObjectAtPositionCommand{
		gameObjs:      objects,
		position:      newPosition,
		objectToPlace: newObject,
		objectAssets:  newObject.GetAssets(),
	}
}

//#endregion

//#region REMOVE OBJECT COMMAND

type removeObjectAtPositionCommand struct {
	gameObjs *GameObjects
	position pixel.Vec
}

func (command *removeObjectAtPositionCommand) execute(waitGroup *sync.WaitGroup) {
	_, index, hit, err := command.gameObjs.getSelectedGameObjAtPosition(command.position)
	if err != nil {
		panic(err)
	}
	if hit {
		*command.gameObjs = command.gameObjs.fastRemoveIndex(index)
	}

	waitGroup.Done()
}

//RemoveObjectAtPosition allows for the removal of a game Object based on Vec location
func (objects *GameObjects) RemoveObjectAtPosition(fromPosition pixel.Vec) ICommand {
	return &removeObjectAtPositionCommand{
		gameObjs: objects,
		position: fromPosition,
	}
}

//#endregion

//#region SELECT OBJECT COMMAND
type selectObjectAtPositionCommand struct {
	gameObjs *GameObjects
	position pixel.Vec
}

func (command *selectObjectAtPositionCommand) execute(waitGroup *sync.WaitGroup) {
	selectedObj, _, hit, err := command.gameObjs.getSelectedGameObjAtPosition(command.position)
	if err != nil {
		panic(err)
	}
	if hit {
		selectedObj.changeState(selected_idle)
	}

	waitGroup.Done()
}

//SelectObjectAtPosition allows for changing the state of a game object based on Vec location to selected
func (objects *GameObjects) SelectObjectAtPosition(fromPosition pixel.Vec) ICommand {
	return &selectObjectAtPositionCommand{
		gameObjs: objects,
		position: fromPosition,
	}
}

//#endregion

//#region MOVE SELECTED OBJECT COMMAND TO POSITION
type moveSelectedObjectToPositionCommand struct {
	gameObjs *GameObjects
	position pixel.Vec
}

func (command *moveSelectedObjectToPositionCommand) execute(waitGroup *sync.WaitGroup) {
	for _, obj := range *command.gameObjs {
		obj.moveToPosition(command.position)
	}
	waitGroup.Done()
}

//MoveSelectedObject allows for directing selected objects to move to a position
func (objects *GameObjects) MoveSelectedToPositionObject(fromPosition pixel.Vec) ICommand {
	return &moveSelectedObjectToPositionCommand{
		gameObjs: objects,
		position: fromPosition,
	}
}

//endregion
