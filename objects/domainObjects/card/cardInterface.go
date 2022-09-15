package card

type ICard interface {
	FlipUp()
	FlipDown()
}

func (card *Card) FlipUp() {
	card.sprite.Set(card.assets.Sheet, card.assets.Anims["up"][0])
}

func (card *Card) FlipDown() {
	card.sprite.Set(card.assets.Sheet, card.assets.Anims["down"][0])
}
