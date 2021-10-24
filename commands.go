package main

import (
	"fmt"
	"sync"

	"github.com/faiface/pixel"
)

//Commands is the map of commands to execute
type Commands map[string]Command

//Command interface is used to execute game commands
type Command interface {
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

type addObjectCommand struct {
	gameObjs       *GameObjects
	objectToPlace  gameObject
	position       pixel.Vec
	animationSheet pixel.Picture
	animations     map[string][]pixel.Rect
	animationKeys  []string
}

func (command *addObjectCommand) execute(waitGroup *sync.WaitGroup) {
	switch command.objectToPlace.(type) {
	case *livingObject:
		{
			*command.gameObjs = command.gameObjs.appendLivingObject(command.animationKeys, command.animations, command.animationSheet, command.position)
		}
	case *gibletObject:
		{
			*command.gameObjs = command.gameObjs.appendGibletObject(command.animationKeys, command.animations, command.animationSheet, command.position)
		}
	}

	waitGroup.Done()
}

//AddObject allows for the addition of a game object
func (objects *GameObjects) AddObject(newObject gameObject, newPosition pixel.Vec) Command {
	return &addObjectCommand{
		gameObjs:       objects,
		position:       newPosition,
		objectToPlace:  newObject,
		animationSheet: newObject.Sheet(),
		animations:     newObject.Animations(),
		animationKeys:  newObject.AnimationKeys(),
	}
}

type removeObjectCommand struct {
	gameObjs *GameObjects
	position pixel.Vec
}

func (command *removeObjectCommand) execute(waitGroup *sync.WaitGroup) {
	selectedObj, index, hit, err := command.gameObjs.getSelectedGameObj(command.position)
	if err != nil {
		panic(err)
	}
	if hit {
		fmt.Println("object id:", selectedObj.getID(), " removed")
		*command.gameObjs = command.gameObjs.fastRemoveIndex(index)
	} else {
		fmt.Println("RemoveObjectCommmand: no object selected")
	}
	hit = false

	waitGroup.Done()
}

//RemoveObject allows for the removal of a game Object based on Vec location
func (objects *GameObjects) RemoveObject(fromPosition pixel.Vec) Command {
	return &removeObjectCommand{
		gameObjs: objects,
		position: fromPosition,
	}
}
