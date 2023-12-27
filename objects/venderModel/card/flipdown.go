package card

import (
	"fmt"

	"github.com/quartermeat/card_game/objects"
)

// FlipDownAction represents the action executed on entering the Off state.
type FlipDownAction struct{}

// Execute Flip down action
func (a *FlipDownAction) Execute(gameObj objects.IGameObject) objects.EventType {

	gameObj.Sprite().Set(gameObj.GetAssets().GetSheet(), gameObj.GetAssets().GetImages()[CARD_BACK])
	fmt.Println("Flip down action")
	return objects.NoOp
}
