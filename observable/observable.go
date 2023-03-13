// Package observable defines the observable pattern in Go
package observable

// ObservableState represents the state of the observable object that observers can observe
type ObservableState struct {
	// TODO: need to define the observables state here
}

// Observable is the object that observers can observe and listen to for changes
type Observable struct {
	ObservableState ObservableState // The state of the observable object
	Subscribers     []IObserver     // The list of observers that are subscribed to the observable
}

// IObserver is the interface for all observers that listen for changes to the observable object
type IObserver interface {
	Notify(appState ObservableState) // The method that observers implement to receive notification when the observable changes
}

// Attach adds an observer to the list of subscribers of the observable object
func (g *Observable) Attach(observer IObserver) {
	g.Subscribers = append(g.Subscribers, observer)
}

// Detach removes an observer from the list of subscribers of the observable object
func (g *Observable) Detach(observer IObserver) {
	for i, obs := range g.Subscribers {
		if obs == observer {
			g.Subscribers = append(g.Subscribers[:i], g.Subscribers[i+1:]...)
			return
		}
	}
}

// NotifyObservers sends a notification to all subscribed observers that the observable has changed
func (g *Observable) NotifyObservers() {
	for _, observer := range g.Subscribers {
		observer.Notify(g.ObservableState)
	}
}

// SetObservableState sets the state of the observable object and notifies all observers that the state has changed
func (g *Observable) SetObservableState(AppState ObservableState) {
	g.ObservableState = AppState
	g.NotifyObservers()
}
