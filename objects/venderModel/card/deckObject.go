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

type Deck struct {
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
func (deck *Deck) ObjectName() string {
	return "Deck"
}

func (deck *Deck) GetFSM() *objects.StateMachine {
	return deck.stateMachine
}

func (deck *Deck) Sprite() *pixel.Sprite {
	return deck.sprite
}

func (deck *Deck) GetAssets() assets.IObjectAsset {
	return deck.asset
}

func (deck *Deck) Selectable() bool {
	return true
}

func (deck *Deck) GetID() int {
	return deck.id
}

func (deck *Deck) Update(dt float64, gameObjects objects.GameObjects, waitGroup *sync.WaitGroup) {
	deck.counter += dt
	// interval := int(math.Floor(card.counter / card.rate))
	//dummy object, with no updates atm
	waitGroup.Done()
}

func (deck *Deck) Draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	for _, card := range deck.cards {
		waitGroup.Add(1)
		card.Draw(win, drawHitBox, waitGroup)
	}

	if drawHitBox {
		imd := imdraw.New(nil)
		imd.Color = pixel.RGB(0, 255, 0)
		imd.Push(deck.GetHitBox().Min, deck.GetHitBox().Max)
		imd.Rectangle(1)
		imd.Draw(win)
	}
	waitGroup.Done()
}

func (deck *Deck) SetHitBox() {
	topRight := pixel.V(deck.position.X-float64(deck.width/2), deck.position.Y-float64(deck.height/2))
	bottomLeft := pixel.V(deck.position.X+float64(deck.width/2), deck.position.Y+float64(deck.width/2))
	deck.hitBox = pixel.R(topRight.X, topRight.Y, bottomLeft.X, bottomLeft.Y)
}

func (deck *Deck) GetHitBox() pixel.Rect {
	return deck.hitBox
}

func (deck *Deck) MoveToPosition(position pixel.Vec) {
	deck.position = position
	deck.matrix = pixel.IM.Moved(position)
	deck.SetHitBox()
}

func (deck *Deck) Shuffle() {
	rand.Shuffle(len(deck.cards), func(i, j int) {
		deck.cards[i], deck.cards[j] = deck.cards[j], deck.cards[i]
	})
}

func (deck *Deck) GetPosition() pixel.Vec{
	return deck.position
}

func (deck *Deck) Deal() ICard {
	if len(deck.cards) == 0 {
		return nil
	}

	card := deck.cards[0]
	deck.cards = deck.cards[1:]

	return card
}

func (deck *Deck) AddCard(card ICard) {
	deck.cards = append(deck.cards, card)
}

func (deck *Deck) PullCard() ICard {
	if len(deck.cards) == 0 {
		return nil
	}

	card := deck.cards[len(deck.cards)-1]
	deck.cards = deck.cards[:len(deck.cards)-1]

	return card
}

func (deck *Deck) GetObservable() *observable.Observable {
	return deck.observable
}

func newDeckFSM() *objects.StateMachine {
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

// NewDeck creates a new deck object containing a set number of card objects
func NewDeckObject(assets assets.ObjectAssets, numCards int, card_type string, position pixel.Vec) Deck {
	deck := Deck{
		id:		 	objects.NextID,
		stateMachine: newDeckFSM(),
		currentState: Operational,
		cards:      make([]ICard, 0, numCards),
		position:   position,
		matrix:     pixel.IM.Moved(position),
		observable: observable.NewObservable(),
		rate:	   1.0,
		dir: 	  0.0,
		vel: 	 pixel.V(0, 0),
	}

	temp_position := deck.position
	
	for i := 0; i < numCards; i++ {
		if(i % 13 == 0)	{
			temp_position.X += 2
			temp_position.Y += 2
		}		
		card := NewCardObject(assets, temp_position, card_type, Hidden)
		deck.cards = append(deck.cards, &card)
		if(i == numCards -1)	{
			deck.width = card.front_sprite.Frame().Max.X - card.front_sprite.Frame().Min.X
			deck.height = card.front_sprite.Frame().Max.Y - card.front_sprite.Frame().Min.Y
		}
	}

	deck.SetHitBox()
	objects.NextID++

	return deck
}
