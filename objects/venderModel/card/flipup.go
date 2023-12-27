package card

import (
	"fmt"

	"github.com/quartermeat/card_game/objects"
)

// OnAction represents the action executed on entering the On state.
type FlipUpAction struct{}

// Execute Flip Up action
func (a *FlipUpAction) Execute(gameObj objects.IGameObject) objects.EventType {
	card := gameObj.(*Card)
	card.currentState = Up
	fmt.Println("Flip up action")
	return objects.NoOp
}
