package card

import (
	"fmt"

	"github.com/quartermeat/card_game/objects"
)

// OnAction represents the action executed on entering the On state.
type FlipUpAction struct{}

// Execute Flip Up action
func (a *FlipUpAction) Execute(gameObj objects.IGameObject) objects.EventType {
	asset := gameObj.GetAssets()		
	gameObj.Sprite().Set(asset.GetSheet(), asset.GetImages()[EVEN_MORE_ZOMBIES])
	fmt.Println("Flip up action")
	return objects.NoOp
}
