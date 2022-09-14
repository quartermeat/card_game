package objects

import (
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/quartermeat/card_game/assets"
)

// can do what both a game Object and an event ctx can do
type IDomainObject interface {
	//IGameObject
	ObjectName() string
	GetFSM() *StateMachine
	Sprite() *pixel.Sprite
	GetAssets() assets.ObjectAsset
	GetID() int
	SetHitBox()
	GetHitBox() pixel.Rect
	Update(dt float64, gameObjects GameObjects, waitGroup *sync.WaitGroup)
	Draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup)
	MoveToPosition(position pixel.Vec)
	SendCtx(object EventContext)
	GetChan() chan EventContext
}

// GetCard returns the card from the game object slice based on the index
// func (gameObjs GameObjects) GetCard(idx int) (card.ICard, error) {

// }
