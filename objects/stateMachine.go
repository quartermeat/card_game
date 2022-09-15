// credit to https://venilnoronha.io/a-simple-state-machine-framework-in-go for FSM framework here
package objects

import (
	"errors"
	"sync"
)

type EventContext struct {
	commandsToObjectChan chan EventContext
}

func (ctx EventContext) GetChan() chan EventContext {
	return ctx.commandsToObjectChan
}

// SendCtx takes a IEventContext, and
// it will asynchronously send it
func (ctx EventContext) SendCtx(object EventContext) {
	objectChannel := EventContext{
		commandsToObjectChan: object.GetChan(),
	}
	select {
	case object.GetChan() <- objectChannel:
		{
			//don't do anything
		}
	default:
		{
			// don't do anything
		}
	}
}

// ErrEventRejected is the error returned when the state machine cannot process
// an event in the state that it is in.
var ErrEventRejected = errors.New("event rejected")

const (
	// Default represents the default state of the system.
	Default StateType = ""

	// NoOp represents a no-op event.
	NoOp EventType = "NoOp"
)

// StateType represents an extensible state type in the state machine.
type StateType string

// EventType represents an extensible event type in the state machine.
type EventType string

// EventContext represents the context to be passed to the action implementation.
// I imagine this could grow with the number of types of objects
type IEventContext interface {
	SendCtx(object EventContext)
	GetChan() chan EventContext
}

// Action represents the action to be executed in a given state.
type Action interface {
	Execute(gameObj IGameObject) EventType
}

// Events represents a mapping of events and states.
type Events map[EventType]StateType

// State binds a state with an action and a set of events it can handle.
type State struct {
	Action Action
	Events Events
}

// States represents a mapping of states and their implementations.
type States map[StateType]State

// StateMachine represents the state machine.
type StateMachine struct {
	// Previous represents the previous state.
	Previous StateType

	// Current represents the current state.
	Current StateType

	// States holds the configuration of states and events handled by the state machine.
	States States

	// mutex ensures that only 1 event is processed by the state machine at any given time.
	mutex sync.Mutex
}

// getNextState returns the next state for the event given the machine's current
// state, or an error if the event can't be handled in the given state.
func (s *StateMachine) getNextState(event EventType) (StateType, error) {
	if state, ok := s.States[s.Current]; ok {
		if state.Events != nil {
			if next, ok := state.Events[event]; ok {
				return next, nil
			}
		}
	}
	return Default, ErrEventRejected
}

// SendEvent sends an event to the state machine.
func (s *StateMachine) SendEvent(event EventType, obj IGameObject) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for {
		// Determine the next state for the event given the machine's current state.
		nextState, err := s.getNextState(event)
		if err != nil {
			return ErrEventRejected
		}

		// Identify the state definition for the next state.
		state, ok := s.States[nextState]
		if !ok || state.Action == nil {
			// configuration error
		}

		// Transition over to the next state.
		s.Previous = s.Current
		s.Current = nextState

		// Execute the next state's action and loop over again if the event returned
		// is not a no-op.
		nextEvent := state.Action.Execute(obj)
		if nextEvent == NoOp {
			return nil
		}
		event = nextEvent
	}
}