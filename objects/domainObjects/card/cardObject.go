// Package 'card' handles the implementation of the card game ojbect
package card

import (
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/quartermeat/card_game/assets"
	"github.com/quartermeat/card_game/objects"
)

// States and Events
const (
	Down objects.StateType = "Down"
	Up   objects.StateType = "Up"

	FlipDown objects.EventType = "FlipDown"
	FlipUp   objects.EventType = "FlipUp"
)

type Card struct {
	stateMachine *objects.StateMachine
	id           int
	assets       assets.ObjectAsset
	sprite       *pixel.Sprite
	rate         float64
	counter      float64
	dir          float64
	vel          pixel.Vec
	hitBox       pixel.Rect
	position     pixel.Vec
	matrix       pixel.Matrix
}

// ObjectName is the string identifier for the object
func (card *Card) ObjectName() string {
	return "Card"
}

func (card *Card) Sprite() *pixel.Sprite {
	return card.sprite
}

func (card *Card) GetAssets() assets.ObjectAsset {
	return card.assets
}

func (card *Card) GetID() int {
	return card.id
}

// SetHitBox sets the hit box for card
// put this in a more efficient location because it won't change per object
func (card *Card) SetHitBox() {
	width := card.sprite.Frame().Max.X - card.sprite.Frame().Min.X
	height := card.sprite.Frame().Max.Y - card.sprite.Frame().Min.Y
	topRight := pixel.V(card.position.X-(width/2), card.position.Y-(height/2))
	bottomLeft := pixel.V(card.position.X+(width/2), card.position.Y+(width/2))
	card.hitBox = pixel.R(topRight.X, topRight.Y, bottomLeft.X, bottomLeft.Y)
}

func (card *Card) GetHitBox() pixel.Rect {
	return card.hitBox
}

func (card *Card) Update(dt float64, gameObjects objects.GameObjects, waitGroup *sync.WaitGroup) {
	card.counter += dt
	// interval := int(math.Floor(card.counter / card.rate))
	//dummy object, with no updates atm
	waitGroup.Done()
}

func (card *Card) Draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	card.sprite.Draw(win, card.matrix)

	if drawHitBox {
		imd := imdraw.New(nil)
		imd.Color = pixel.RGB(0, 255, 0)
		imd.Push(card.hitBox.Min, card.hitBox.Max)
		imd.Rectangle(1)
		imd.Draw(win)
	}
	waitGroup.Done()
}

func (card *Card) MoveToPosition(position pixel.Vec) {
	panic("not implemented") // TODO: Implement
}

func (card *Card) GetFSM() *objects.StateMachine {
	return card.stateMachine
}

func newCardFSM() *objects.StateMachine {
	return &objects.StateMachine{
		States: objects.States{
			objects.Default: objects.State{
				Events: objects.Events{
					FlipDown: Down,
				},
			},
			Down: objects.State{
				Action: &FlipUpAction{},
				Events: objects.Events{
					FlipUp: Up,
				},
			},
			Up: objects.State{
				Action: &FlipDownAction{},
				Events: objects.Events{
					FlipDown: Down,
				},
			},
		},
	}
}

// NewCardObject creates a new card game object
func NewCardObject(objectAsset assets.ObjectAsset, position pixel.Vec) Card {
	downAnimationKey := objectAsset.AnimKeys[0]
	downAnimationFrame := 0

	newCard := Card{
		id:           objects.NextID,
		stateMachine: newCardFSM(),
		assets:       objectAsset,
		sprite:       pixel.NewSprite(objectAsset.Sheet, objectAsset.Anims[downAnimationKey][downAnimationFrame]),
		rate:         1.0,
		dir:          0.0,
		position:     position,
		vel:          pixel.V(0, 0),
		matrix:       pixel.IM.Moved(position),
	}
	newCard.SetHitBox()
	objects.NextID++
	return newCard
}
