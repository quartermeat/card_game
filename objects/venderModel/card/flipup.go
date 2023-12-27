package card

import (
	"fmt"

	"github.com/quartermeat/card_game/objects"
)

// OnAction represents the action executed on entering the On state.
type FlipUpAction struct{}

// Execute Flip Up action
func (a *FlipUpAction) Execute(gameObj objects.IGameObject) objects.EventType {
	
	assets := gameObj.GetAssets()
	sheet := assets.GetSheet()
	images := assets.GetImages()
	
	fmt.Printf("setting sprite: assets:%s/n", assets.GetDescription())

	for image := range images {
		fmt.Printf("image: %s\n", image)
	}

	gameObj.Sprite().Set(sheet, images[AMMO_BOX])
	fmt.Println("Flip up action")
	return objects.NoOp
}
