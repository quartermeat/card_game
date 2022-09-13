package card

import (
	"fmt"

	"github.com/quartermeat/card_game/objects"
)

// OnAction represents the action executed on entering the On state.
type FlipUpAction struct{}

// Execute Flip Up action
func (a *FlipUpAction) Execute(eventCtx objects.EventContext) objects.EventType {
	fmt.Println("Flip up action")
	return objects.NoOp
}
