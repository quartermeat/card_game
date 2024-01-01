package card

import (
	"fmt"

	"github.com/quartermeat/card_game/objects"
)

// PullAction represents the pull action for a card.
type PullAction struct{}

// Execute executes the pull action.
func (pa *PullAction) Execute(gameObj objects.IGameObject) objects.EventType{
	fmt.Printf("Pulling card\n")
	// Add your pull logic here
	deck := gameObj.(IDeck)

	deck.PullCard()

	return objects.NoOp
}

