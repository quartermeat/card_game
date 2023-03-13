// Package 'Observable' is a go package to provide structures that have state changes observable to subscribers
package observable

//ObservableState is used by ipc package for testing communication between 'threads'
type ObservableState struct {
	// TODO: need to define app state here...
}

type Observable struct {
	ObservableState ObservableState
	Observers       []IObserver
}

type IObserver interface {
	Notify(appState ObservableState)
}

func (g *Observable) Attach(observer IObserver) {
	g.Observers = append(g.Observers, observer)
}

func (g *App) Detach(observer IObserver) {
	for i, obs := range g.Observers {
		if obs == observer {
			g.Observers = append(g.Observers[:i], g.Observers[i+1:]...)
			return
		}
	}
}

func (g *App) NotifyObservers() {
	for _, observer := range g.Observers {
		observer.Notify(g.AppState)
	}
}

func (g *App) SetAppState(AppState AppState) {
	g.AppState = AppState
	g.NotifyObservers()
}
