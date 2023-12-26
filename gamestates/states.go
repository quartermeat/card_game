package gamestates

type State int

const (
	Init State = iota
	Ready
)

type StateManager struct {
	currentState State
}

func NewStateManager() *StateManager {
	return &StateManager{
		currentState: Init,
	}
}

func (sm *StateManager) GetCurrentState() State {
	return sm.currentState
}

func (sm *StateManager) SetCurrentState(newState State) {
	sm.currentState = newState
}
