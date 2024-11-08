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

type Hand struct {
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
func (hand *Hand) ObjectName() string {
	return "Hand"
}

func (hand *Hand) Selectable() bool {
	return true
}

func (hand *Hand) GetFSM() *objects.StateMachine {
	return hand.stateMachine
}

func (hand *Hand) Sprite() *pixel.Sprite {
	return hand.sprite
}

func (hand *Hand) GetAssets() assets.IObjectAsset {
	return hand.asset
}

func (hand *Hand) GetID() int {
	return hand.id
}

func (hand *Hand) Update(dt float64, gameObjects objects.GameObjects, waitGroup *sync.WaitGroup) {
	hand.counter += dt
	// interval := int(math.Floor(card.counter / card.rate))
	//dummy object, with no updates atm
	waitGroup.Done()
}

// func (hand *Hand) Draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
//     // Calculate the spread of cards in radians, e.g., PI/6 for a 30° spread
//     spread := math.Pi / 6 // Adjust this value to control the spread of the fan
//     startAngle := -spread * float64(len(hand.cards)-1) / 2 // Starting angle for the first card

//     for i, card := range hand.cards {
//         waitGroup.Add(1)

//         // Calculate the angle for the current card
//         angle := startAngle + spread*float64(i)
//         // Calculate the position offset based on the angle
//         // Adjust the radius to control how far out the cards spread from the center
//         radius := 50.0 // This controls how spread out the cards are
//         offsetX := math.Cos(angle) * radius
//         offsetY := math.Sin(angle) * radius

//         // Update card's position based on the calculated offset
//         cardPosition := pixel.V(hand.position.X+offsetX, hand.position.Y+offsetY)
//         // Update card's rotation based on the angle
//         cardMatrix := pixel.IM.Rotated(pixel.ZV, angle).Moved(cardPosition)

//         // Use the card's updated matrix for drawing
//         cardSprite := card.Sprite()
//         cardSprite.SetMatrix(cardMatrix)
//         cardSprite.Draw(win)

//         if drawHitBox {
//             imd := imdraw.New(nil)
//             imd.Color = pixel.RGB(0, 255, 0)
//             // Adjust this to draw the hitbox around the rotated card if necessary
//             imd.Push(cardPosition, cardPosition.Add(pixel.V(10, 10))) // Example; adjust as needed
//             imd.Rectangle(1)
//             imd.Draw(win)
//         }

//         waitGroup.Done()
//     }

//     if drawHitBox {
//         imd := imdraw.New(nil)
//         imd.Color = pixel.RGB(0, 255, 0)
//         imd.Push(hand.GetHitBox().Min, hand.GetHitBox().Max)
//         imd.Rectangle(1)
//         imd.Draw(win)
//     }
//     waitGroup.Done()
// }

func (hand *Hand) Draw(win *pixelgl.Window, drawHitBox bool, waitGroup *sync.WaitGroup) {
	for _, card := range hand.cards {
		waitGroup.Add(1)
		card.Draw(win, drawHitBox, waitGroup)
	}

	if drawHitBox {
		imd := imdraw.New(nil)
		imd.Color = pixel.RGB(0, 255, 0)
		imd.Push(hand.GetHitBox().Min, hand.GetHitBox().Max)
		imd.Rectangle(1)
		imd.Draw(win)
	}
	waitGroup.Done()
}

func (hand *Hand) SetHitBox() {
	topRight := pixel.V(hand.position.X-float64(hand.width/2), hand.position.Y-float64(hand.height/2))
	bottomLeft := pixel.V(hand.position.X+float64(hand.width/2), hand.position.Y+float64(hand.width/2))
	hand.hitBox = pixel.R(topRight.X, topRight.Y, bottomLeft.X, bottomLeft.Y)
}

func (hand *Hand) GetHitBox() pixel.Rect {
	return hand.hitBox
}

func (hand *Hand) MoveToPosition(position pixel.Vec) {
	hand.position = position
	hand.matrix = pixel.IM.Moved(position)
	hand.SetHitBox()
}

func (hand *Hand) Shuffle() {
	rand.Shuffle(len(hand.cards), func(i, j int) {
		hand.cards[i], hand.cards[j] = hand.cards[j], hand.cards[i]
	})
}

func (hand *Hand) Deal() ICard {
	if len(hand.cards) == 0 {
		return nil
	}

	card := hand.cards[0]
	hand.cards = hand.cards[1:]

	return card
}

func (hand *Hand) GetPosition() pixel.Vec{
	return hand.position
}

func (hand *Hand) AddCard(card ICard) {
	hand.cards = append(hand.cards, card)
}

func (hand *Hand) PullCard() ICard {
	if len(hand.cards) == 0 {
		return nil
	}

	card := hand.cards[0]
	hand.cards = hand.cards[1:]

	return card
}

func (hand *Hand) GetObservable() *observable.Observable {
	return hand.observable
}

func newHandFSM() *objects.StateMachine {
	return &objects.StateMachine{
		States: objects.States{
			objects.Default: objects.State{
				Events: objects.Events{
					Play: Operational,
				},
			},
			Operational: objects.State{
				Action: &PlayAction{},
				Events: objects.Events{
					Play: Operational,
				},
			},
			Empty: objects.State{
				Events: objects.Events{	
				},
			},
		},
	}
}

// NewHand creates a new hand object containing a set number of card objects
func NewHandObject(assets assets.ObjectAssets, position pixel.Vec) Hand {
	numCards := 5
	
	hand := Hand{
		id:		 	objects.NextID,
		stateMachine: newHandFSM(),
		currentState: Operational,
		cards:      make([]ICard, 0, numCards),
		position:   position,
		matrix:     pixel.IM.Moved(position),
		observable: observable.NewObservable(),
		rate:	   1.0,
		dir: 	  0.0,
		vel: 	 pixel.V(0, 0),
	}

	//need to implement to setup a default hand with specific cards per dominion rules

	temp_position := hand.position
	
	for i := 0; i < numCards; i++ {
		temp_position.X += 2
		temp_position.Y += 2
		card := NewCardObject(assets, temp_position, "zombies", Hidden)
		hand.cards = append(hand.cards, &card)
		if(i == numCards -1)	{
			hand.width = card.front_sprite.Frame().Max.X - card.front_sprite.Frame().Min.X
			hand.height = card.front_sprite.Frame().Max.Y - card.front_sprite.Frame().Min.Y
		}
	}

	hand.SetHitBox()
	objects.NextID++

	return hand
}
