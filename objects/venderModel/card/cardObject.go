// Package 'card' handles the implementation of the card game ojbect
package card

import (
	"sync"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/quartermeat/card_game/assets"
	"github.com/quartermeat/card_game/objects"
	"github.com/quartermeat/card_game/observable"
)

// States and Events
const (
	Down objects.StateType = "Down"
	Up   objects.StateType = "Up"

	Flip objects.EventType = "Flip"
)

type Card struct {
	stateMachine *objects.StateMachine
	currentState objects.StateType
	id           int
	asset        assets.ObjectImageAsset
	front_sprite *pixel.Sprite
	back_sprite	 *pixel.Sprite
	rate         float64
	counter      float64
	dir          float64
	vel          pixel.Vec
	hitBox       pixel.Rect
	position     pixel.Vec
	matrix       pixel.Matrix
	observable   *observable.Observable
}

// ObjectName is the string identifier for the object
func (card *Card) ObjectName() string {
	return "Card"
}

func (card *Card) Sprite() *pixel.Sprite {
	return card.front_sprite
}

func (card *Card) GetAssets() assets.IObjectAsset {
	return card.asset
}

func (card *Card) GetID() int {
	return card.id
}

// SetHitBox sets the hit box for card
// TODO: put this in a more efficient location because it won't change per object
func (card *Card) SetHitBox() {
	width := card.front_sprite.Frame().Max.X - card.front_sprite.Frame().Min.X
	height := card.front_sprite.Frame().Max.Y - card.front_sprite.Frame().Min.Y
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

func (card *Card) GetObservable() *observable.Observable {
	return card.observable
}

func (card *Card) Draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	if(card.currentState == Down) {
		card.back_sprite.Draw(win, card.matrix)
	} else {
		card.front_sprite.Draw(win, card.matrix)
	}
	
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
				Action: &FlipAction{},
				Events: objects.Events{
					Flip: Up,
				},
			},
			Down: objects.State{
				Action: &FlipAction{},
				Events: objects.Events{
					Flip: Up,
				},
			},
			Up: objects.State{
				Action: &FlipAction{},
				Events: objects.Events{
					Flip: Down,
				},
			},
		},
	}
}

func (card *Card) GetBackSprite() *pixel.Sprite {
	return card.back_sprite
}

func (card* Card) GetMatrix() pixel.Matrix {
	return card.matrix
}

// NewCardObject creates a new card game object
func NewCardObject(objectAssets assets.ObjectAssets, position pixel.Vec, card_name string) Card {
	objectAsset := objectAssets.GetImage(card_name)
	objAsset := objectAsset.(assets.ObjectImageAsset)

	backAsset := objectAssets.GetImage(CARD_BACK)
	backObjAsset := backAsset.(assets.ObjectImageAsset)

	newCard := Card{
		id:           objects.NextID,
		currentState: Down,
		stateMachine: newCardFSM(),
		asset:       objectAsset.(assets.ObjectImageAsset),
		front_sprite: pixel.NewSprite(objAsset.Sheet, objAsset.GetImages()[card_name]),
		back_sprite:  pixel.NewSprite(backObjAsset.Sheet, backObjAsset.GetImages()[CARD_BACK]),
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
