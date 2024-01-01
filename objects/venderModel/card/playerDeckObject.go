package card

import (
	"math/rand"
	"sync"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/quartermeat/card_game/assets"
	"github.com/quartermeat/card_game/objects"
	"github.com/quartermeat/card_game/observable"
)

type PlayerDeck struct {
	cards      	 []ICard
	asset 		 assets.ObjectImageAsset
	position   	 pixel.Vec
	hitBox     	 pixel.Rect
	matrix     	 pixel.Matrix
	observable 	 *observable.Observable
	stateMachine *objects.StateMachine
	currentState objects.StateType
	id 			 int
	counter      float64
	vel          pixel.Vec
	dir          float64
	rate         float64
	sprite	     *pixel.Sprite
	height		 float64
	width		 float64
}

// ObjectName is the string identifier for the object
func (playerDeck *PlayerDeck) ObjectName() string {
	return "PlayerDeck"
}

func (playerDeck *PlayerDeck) GetFSM() *objects.StateMachine {
	return playerDeck.stateMachine
}

func (playerDeck *PlayerDeck) Sprite() *pixel.Sprite {
	return playerDeck.sprite
}

func (playerDeck *PlayerDeck) GetAssets() assets.IObjectAsset {
	return playerDeck.asset
}

func (playerDeck *PlayerDeck) Selectable() bool {
	return true
}

func (playerDeck *PlayerDeck) GetID() int {
	return playerDeck.id
}

func (playerDeck *PlayerDeck) Update(dt float64, gameObjects objects.GameObjects, waitGroup *sync.WaitGroup) {
	playerDeck.counter += dt
	// interval := int(math.Floor(card.counter / card.rate))
	//dummy object, with no updates atm
	waitGroup.Done()
}

func (playerDeck *PlayerDeck) Draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	for _, card := range playerDeck.cards {
		waitGroup.Add(1)
		card.Draw(win, drawHitBox, waitGroup)
	}

	if drawHitBox {
		imd := imdraw.New(nil)
		imd.Color = pixel.RGB(0, 255, 0)
		imd.Push(playerDeck.GetHitBox().Min, playerDeck.GetHitBox().Max)
		imd.Rectangle(1)
		imd.Draw(win)
	}
	waitGroup.Done()
}

func (playerDeck *PlayerDeck) SetHitBox() {
	topRight := pixel.V(playerDeck.position.X-float64(playerDeck.width/2), playerDeck.position.Y-float64(playerDeck.height/2))
	bottomLeft := pixel.V(playerDeck.position.X+float64(playerDeck.width/2), playerDeck.position.Y+float64(playerDeck.width/2))
	playerDeck.hitBox = pixel.R(topRight.X, topRight.Y, bottomLeft.X, bottomLeft.Y)
}

func (playerDeck *PlayerDeck) GetHitBox() pixel.Rect {
	return playerDeck.hitBox
}

func (playerDeck *PlayerDeck) MoveToPosition(position pixel.Vec) {
	playerDeck.position = position
	playerDeck.matrix = pixel.IM.Moved(position)
	playerDeck.SetHitBox()
}

func (playerDeck *PlayerDeck) Shuffle() {
	rand.Shuffle(len(playerDeck.cards), func(i, j int) {
		playerDeck.cards[i], playerDeck.cards[j] = playerDeck.cards[j], playerDeck.cards[i]
	})
}

func (playerDeck *PlayerDeck) Deal() ICard {
	if len(playerDeck.cards) == 0 {
		return nil
	}

	card := playerDeck.cards[0]
	playerDeck.cards = playerDeck.cards[1:]

	return card
}

func (playerDeck *PlayerDeck) AddCard(card ICard) {
	playerDeck.cards = append(playerDeck.cards, card)
}

func (playerDeck *PlayerDeck) PullCard() ICard {
	if len(playerDeck.cards) == 0 {
		return nil
	}

	card := playerDeck.cards[len(playerDeck.cards)-1]
	playerDeck.cards = playerDeck.cards[:len(playerDeck.cards)-1]

	return card
}

func (playerDeck *PlayerDeck) GetObservable() *observable.Observable {
	return playerDeck.observable
}

func newPlayerDeckFSM() *objects.StateMachine {
	return &objects.StateMachine{
		States: objects.States{
			objects.Default: objects.State{
				Action: &PullAction{},
				Events: objects.Events{
					Pull: Operational,
				},
			},
			Operational: objects.State{
				Action: &PullAction{},
				Events: objects.Events{
					Pull: Operational,
				},
			},
			Empty: objects.State{
				Events: objects.Events{	
				},
			},
		},
	}
}

// NewPlayerDeckObject creates a new playerDeck object containing a set number of card objects
func NewPlayerDeckObject(assets assets.ObjectAssets, position pixel.Vec) PlayerDeck {
	playerDeck := PlayerDeck{
		id:		 	objects.NextID,
		stateMachine: newPlayerDeckFSM(),
		currentState: Operational,
		cards:      make([]ICard, 0, 10),
		position:   position,
		matrix:     pixel.IM.Moved(position),
		observable: observable.NewObservable(),
		rate:	   1.0,
		dir: 	  0.0,
		vel: 	 pixel.V(0, 0),
	}

	temp_position := playerDeck.position
	var card_type string

	for i := 0; i < 10; i++ {
		if(i % 13 == 0)	{
			temp_position.X += 2
			temp_position.Y += 2
		}
		if(i < 3){
			card_type = "zombies"
		}else{
			card_type = "bullet"
		}
		card := NewCardObject(assets, temp_position, card_type, Down)
		playerDeck.cards = append(playerDeck.cards, &card)
		if(i == 9)	{
			playerDeck.width = card.front_sprite.Frame().Max.X - card.front_sprite.Frame().Min.X
			playerDeck.height = card.front_sprite.Frame().Max.Y - card.front_sprite.Frame().Min.Y
		}
	}

	playerDeck.Shuffle()

	playerDeck.SetHitBox()
	objects.NextID++

	return playerDeck
}
