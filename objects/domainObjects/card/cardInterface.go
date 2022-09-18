package card

type ICard interface {
	FlipUp()
	FlipDown()
}

func (card *Card) FlipUp() {
	card.sprite.Set(card.assets.Sheet, card.assets.GetImages()["sacrifice"])
}

func (card *Card) FlipDown() {
	card.sprite.Set(card.assets.Sheet, card.assets.GetImages()["shotgun"])
}
