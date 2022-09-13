package card

import (
	"fmt"

	"github.com/quartermeat/card_game/objects"
)

// FlipDownAction represents the action executed on entering the Off state.
type FlipDownAction struct{}

// Execute Flip down action
func (a *FlipDownAction) Execute(eventCtx objects.EventContext) objects.EventType {

	fmt.Println("Flip down action")
	return objects.NoOp
}
