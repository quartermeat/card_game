package input

import (
	"fmt"
	"math/rand"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/quartermeat/card_game/assets"
	"github.com/quartermeat/card_game/objects"
	"github.com/quartermeat/card_game/objects/venderModel/card"
)

const (
	CursorDesription = "hand"
)

var (
	// Selectable Objects
	Card = "Card"
	Deck = "Deck"
	PlayerDeck = "PlayerDeck"
)

var (
	// Object Events
	Flip = card.Flip
	Pull = card.Pull
)

// seems I can't stack commands, so InitGame has to happen in stages
var (
	Stage = 0
	kingdom_card_bag = []string{"ham_radio", 
								"survivors",
								"1_in_the_chamber",
								"ammo_box",
								"barricade",
								"courage",
								"cunning",
								"decoy",
								"hide",
								"higher_ground",
								"hollow_points",
								"maverick",
								"molotov_cocktail",
								"quick_escape",
								"recon",
								"regroup",
								"reload",
								"restock",
								"sacrifice",
								"scavenger",
								"shotgun",
								"sidekick",
								"stick_together",
								"zombie_swarm",
								"tactics",
								"weapons_cache",
								}
)

func getRandomKingdomCard() string {
	index := rand.Intn(len(kingdom_card_bag))
	card := kingdom_card_bag[index]
	kingdom_card_bag = append(kingdom_card_bag[:index], kingdom_card_bag[index+1:]...)	
	return card
}

// setup game board, create objects, and setup input
func InitGame(win *pixelgl.Window, cam *pixel.Matrix, gameCommands Commands, gameObjs *objects.GameObjects, objectAssets assets.ObjectAssets) bool {
	
	// setup a deck of cards positioned on the wooden background
	// start with a top row of 6 decks, treaure and estates
	//top left corner: x:-900, y:1200
	startx := -900.0
	starty := 1100.0
	rowy := 700.0
	var objectToPlace objects.IGameObject
	victory_point_deck_size := 8
	kingdom_card_deck_size := 8

	switch(Stage){
	case 0:{
		//first deck is the bullet deck
		bullet_deck_size := 80
		location := pixel.Vec{X: startx, Y: starty}
		objectToPlace := card.NewDeckObject(objectAssets, bullet_deck_size, "bullet", location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)
		Stage = 1
		return false
	}
	case 1:{
		// to the right of the bullet deck is the slug deck
		slug_deck_size := 70
		location := pixel.Vec{X: startx + 250, Y: starty}
		objectToPlace := card.NewDeckObject(objectAssets, slug_deck_size, "slug", location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 2
		return false
	}
	case 2:{
		// to the right of the slug deck is the shell deck
		shell_deck_size := 48
		location := pixel.Vec{X: startx + 500, Y: starty}
		objectToPlace := card.NewDeckObject(objectAssets, shell_deck_size, "shells", location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 3
		return false
	}
	case 3:{
		// to the right of the shell deck is the start of the 'estate' cards, first with zombies
		location := pixel.Vec{X: startx + 800, Y: starty}
		objectToPlace := card.NewDeckObject(objectAssets, victory_point_deck_size, "zombies", location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 4
		return false
	}
	case 4:{
		// to the right of the zombie deck is the more_zombies deck
		location := pixel.Vec{X: startx + 1050, Y: starty}
		objectToPlace := card.NewDeckObject(objectAssets, victory_point_deck_size, "more_zombies", location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 5
		return false
	}
	case 5:{
		// to the right of the more_zombies deck is the even_more_zombies deck
		location := pixel.Vec{X: startx + 1300, Y: starty}
		objectToPlace := card.NewDeckObject(objectAssets, victory_point_deck_size, "even_more_zombies", location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 6
		return false
	}
	case 6:{
		// to the right of the even_more_zombies deck is the infections deck
		infections_deck_size := 10
		location := pixel.Vec{X: startx + 1550, Y: starty}
		objectToPlace := card.NewDeckObject(objectAssets, infections_deck_size, "infection", location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 7
		return false
	}
	case 7:{
		// Kingdom card 1
		
		location := pixel.Vec{X: startx, Y: rowy}
		objectToPlace := card.NewDeckObject(objectAssets, kingdom_card_deck_size, getRandomKingdomCard(), location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 8
		return false
	}
	case 8:{
		// Kingdom card 2
		location := pixel.Vec{X: startx + 250, Y: rowy}
		objectToPlace := card.NewDeckObject(objectAssets, kingdom_card_deck_size, getRandomKingdomCard(), location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 9
		return false
	}
	case 9:{
		// Kingdom card 3
		location := pixel.Vec{X: startx + 500, Y: rowy}
		objectToPlace := card.NewDeckObject(objectAssets, kingdom_card_deck_size, getRandomKingdomCard(), location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 10
		return false
	}
	case 10:{
		// Kingdom card 4
		location := pixel.Vec{X: startx + 750, Y: rowy}
		objectToPlace := card.NewDeckObject(objectAssets, kingdom_card_deck_size, getRandomKingdomCard(), location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 11
		return false
	}
	case 11:{
		// Kingdom card 5
		location := pixel.Vec{X: startx + 1000, Y: rowy}
		objectToPlace := card.NewDeckObject(objectAssets, kingdom_card_deck_size, getRandomKingdomCard(), location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 12
		return false
	}
	case 12:{
		// Kingdom card 6
		location := pixel.Vec{X: startx + 1250, Y: rowy}
		objectToPlace := card.NewDeckObject(objectAssets, kingdom_card_deck_size, getRandomKingdomCard(), location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 13
		return false
	}
	case 13:{
		// Kingdom card 7
		location := pixel.Vec{X: startx + 1500, Y: rowy}
		objectToPlace := card.NewDeckObject(objectAssets, kingdom_card_deck_size, getRandomKingdomCard(), location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 14
		return false
	}
	case 14:{
		// Kingdom card 8
		location := pixel.Vec{X: startx + 1750, Y: rowy}
		objectToPlace := card.NewDeckObject(objectAssets, kingdom_card_deck_size, getRandomKingdomCard(), location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 15
		return false
	}
	case 15:{
		// Kingdom card 9
		location := pixel.Vec{X: startx + 2000, Y: rowy}
		objectToPlace := card.NewDeckObject(objectAssets, kingdom_card_deck_size, getRandomKingdomCard(), location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 16
		return false
	}
	case 16:{
		// Kingdom card 10
		location := pixel.Vec{X: startx + 2250, Y: rowy}
		objectToPlace := card.NewDeckObject(objectAssets, kingdom_card_deck_size, getRandomKingdomCard(), location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 17
		return false
	}
	case 17:{
		// Player Deck setup
		// executing: SelectObjectAtPosition x:-394.317658, y:-295.212168
		location := pixel.Vec{X: -400, Y: -300}
		objectToPlace := card.NewPlayerDeckObject(objectAssets, location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 18
		return false
	}
	case 18:{
		// AI Deck setup
		location := pixel.Vec{X: 2000, Y: -300}
		objectToPlace := card.NewPlayerDeckObject(objectAssets, location)
		gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
		Stage = 19
		return false
	}
	case 19:{
		// Player Hand setup

		// location := pixel.Vec{X: 2000, Y: -300}
		// objectToPlace := card.NewPlayerDeckObject(objectAssets, location)
		// gameCommands[fmt.Sprintf("AddObjectAtPosition: x:%f, y:%f, ObjectType:%s", location.X, location.Y, objectToPlace.ObjectName())] = AddObjectAtPosition(gameObjs, &objectToPlace, location)	
				
		Stage = 20
		return true
	}
	default:{
		fmt.Printf("InitGame: Stage %d is not defined, objectToPlace: %s", Stage, objectToPlace.ObjectName())
	}
	}

	// arrange the kingdom in the middle
	
	return false
}