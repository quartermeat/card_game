package objects

type LivingIdleState struct {
	//move as much state here as we need to make changes to
}

func (state *LivingIdleState) GetState() string {
	return "idle"
}

func (state *LivingIdleState) EnterState() {
	// obj.counter = 0
	// obj.matrix = pixel.IM.Moved(obj.position)

	// case MOVING:
	// 	{
	// 		if livingObj.motives.destinationReached {
	// 			randFloat := float64(rand.Intn(360))
	// 			livingObj.dir = randFloat * (math.Pi / 180)
	// 		}
	// 		livingObj.matrix = livingObj.matrix.Rotated(livingObj.position, livingObj.dir)
	// 	}
	// case SELECTED_IDLE:
	// 	{
	// 		livingObj.matrix = pixel.IM.Moved(livingObj.position)
	// 	}
	// case SELECTED_MOVING:
	// 	{
	// 		livingObj.motives.destinationReached = false
	// 		livingObj.ChangeState(MOVING)
	// 	}
	// }
}

func (state *LivingIdleState) GoMoving(obj *LivingObject) {
	//start moving
}

func (state *LivingIdleState) DoMoving(obj *LivingObject) {
	//start moving
}

func (state *LivingIdleState) GoIdle(obj *LivingObject) {
	//go idle
}

func (state *LivingIdleState) DoIdle(obj *LivingObject, interval int) {
	//update idle animation
	obj.sprite.Set(obj.assets.Sheet, obj.assets.Anims["idle"][interval%len(obj.assets.Anims["idle"])])

	obj.attributes.stamina += obj.counter

	//start moving in a random direction
	// if obj.counter >= obj.attributes.initiative {
	// 	obj.ChangeState(MOVING)
	// }
}
