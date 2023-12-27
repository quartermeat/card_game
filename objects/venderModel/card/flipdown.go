package card

import (
	"fmt"

	"github.com/quartermeat/card_game/objects"
)

// FlipDownAction represents the action executed on entering the Off state.
type FlipAction struct{}

// Execute Flip down action
func (a *FlipAction) Execute(gameObj objects.IGameObject) objects.EventType {
	card := gameObj.(*Card)
	if(card.currentState == Up){
		card.currentState = Down
	}else {
		card.currentState = Up
	}
	fmt.Println("Flip action")
	fmt.Printf("Current State: %s\n", card.currentState)
	return objects.NoOp
}
