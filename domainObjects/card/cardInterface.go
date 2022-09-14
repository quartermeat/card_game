package card

import "github.com/quartermeat/card_game/objects"

type ICard interface {
	FlipUp()
	FlipDown()
	GetChan() chan objects.EventContext
	SendCtx(objects.EventContext)
}

func (card *Card) FlipUp() {
	card.sprite.Set(card.assets.Sheet, card.assets.Anims["up"][0])
}

func (card *Card) FlipDown() {
	card.sprite.Set(card.assets.Sheet, card.assets.Anims["down"][0])
}

func (card *Card) GetChan() chan objects.EventContext {
	return card.commandsToCard
}

func (card *Card) SendCtx(object objects.EventContext) {
	select {
	case object.GetChan() <- <-card.commandsToCard:
		{
			//don't do anything
		}
	default:
		{
			// don't do anything
		}
	}

}
