package card

import (
	"math/rand"
	"sync"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/quartermeat/card_game/assets"
	"github.com/quartermeat/card_game/objects"
	"github.com/quartermeat/card_game/observable"
)

const (
	DeckWidth  = 70
	DeckHeight = 100
)

// States and Events
const (
	Operational objects.StateType = "Operational"
	Empty objects.StateType = "Empty"
	Pull objects.EventType = "Pull"	
)

type Deck struct {
	cards      []ICard
	position   pixel.Vec
	hitBox     pixel.Rect
	matrix     pixel.Matrix
	observable *observable.Observable
	stateMachine *objects.StateMachine
	image string		
}

// ObjectName is the string identifier for the object
func (deck *Deck) ObjectName() string {
	return "Deck"
}

func (deck *Deck) GetAssets() assets.IObjectAsset {
	return nil
}

func (deck *Deck) GetID() int {
	return 0
}

func (deck *Deck) Update(dt float64, gameObjects objects.GameObjects, waitGroup *sync.WaitGroup) {
	// TODO: Implement
	waitGroup.Done()
}

func (deck *Deck) Draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	// Use deck.image to draw the deck
	// draw the same image with each image offset a little to see the one below it, draw as many cards as there are in the deck
	waitGroup.Done()
}

func (deck *Deck) SetHitBox() {
	width := DeckWidth
	height := DeckHeight
	topRight := pixel.V(deck.position.X-float64(width/2), deck.position.Y-float64(height/2))
	bottomLeft := pixel.V(deck.position.X+float64(width/2), deck.position.Y+float64(width/2))
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

func (deck *Deck) GetObservable() *observable.Observable {
	return deck.observable
}

func newDeckFSM() *objects.StateMachine {
	return &objects.StateMachine{
		States: objects.States{
			objects.Default: objects.State{
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
func NewDeck(assets assets.ObjectAssets, numCards int, imagePath string) *Deck {
	deck := &Deck{
		stateMachine: newDeckFSM(),
		cards:      make([]ICard, 0, numCards),
		position:   pixel.V(0, 0),
		matrix:     pixel.IM.Moved(pixel.V(0, 0)),
		observable: observable.NewObservable(),
		image: imagePath,
	}

	for i := 0; i < numCards; i++ {
		card := NewCardObject(nil, pixel.V(0, 0), "")
		deck.cards = append(deck.cards, &card)
	}

	deck.SetHitBox()

	return deck
}
