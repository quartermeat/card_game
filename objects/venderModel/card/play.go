package card

import (
	"fmt"

	"github.com/quartermeat/card_game/objects"
)

// PullAction represents the pull action for a card.
type PlayAction struct{}

// Execute executes the pull action.
func (pa *PlayAction) Execute(gameObj objects.IGameObject) objects.EventType{
	fmt.Printf("Playing card from hand\n")
	// Add play card logic here
	// hand := gameObj.(IHand)

	return objects.NoOp
}

