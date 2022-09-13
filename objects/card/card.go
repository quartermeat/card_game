package card

import (
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/quartermeat/card_game/assets"
	"github.com/quartermeat/card_game/objects"
)

type Card struct {
}

func (card Card) ObjectName() string {
	panic("not implemented") // TODO: Implement
}

func (card Card) Sprite() *pixel.Sprite {
	panic("not implemented") // TODO: Implement
}

func (card Card) GetAssets() assets.ObjectAssets {
	panic("not implemented") // TODO: Implement
}

func (card Card) GetID() int {
	panic("not implemented") // TODO: Implement
}

func (card Card) SetHitBox() {
	panic("not implemented") // TODO: Implement
}

func (card Card) GetHitBox() pixel.Rect {
	panic("not implemented") // TODO: Implement
}

func (card Card) Update(dt float64, gameObjects objects.GameObjects, waitGroup *sync.WaitGroup) {
	panic("not implemented") // TODO: Implement
}

func (card Card) ChangeControlState(newState objects.IControlState) {
	panic("not implemented") // TODO: Implement
}

func (card Card) Draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	panic("not implemented") // TODO: Implement
}

func (card Card) MoveToPosition(position pixel.Vec) {
	panic("not implemented") // TODO: Implement
}
