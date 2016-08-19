package lab

import "sync"

// Lab is a playground where you can define experiments and strategies.
type Lab interface {
	// DefineStrategy defines strategies to aim at visitors.
	DefineStrategy(string, Strategy)
	// Experiment defines an experiment or retrieves one that already exists.
	Experiment(string) Experiment
	// Session starts a new session for the given visitor.
	Session(Visitor) Session
}

// Strategy will tell if the given Params should be shown the experiment.
type Strategy func(Params) bool

// Visitor is something identified by an ID. Represents the visitor of a
// single session so all features are shown or not during the same session.
type Visitor interface {
	// ID returns the ID of the visitor.
	ID() string
}

type lab struct {
	sync.RWMutex
	strategies  map[string]Strategy
	experiments map[string]*experiment
}

// New returns a new lab.
func New() Lab {
	return &lab{
		strategies: make(map[string]Strategy),
	}
}

func (l *lab) Strategy(ID string) (Strategy, bool) {
	l.Lock()
	defer l.Unlock()
	s, ok := l.strategies[ID]
	return s, ok
}

func (l *lab) DefineStrategy(ID string, strategy Strategy) {
	l.Lock()
	defer l.Unlock()
	l.strategies[ID] = strategy
}

func (l *lab) Experiment(ID string) Experiment {
	return l.experiment(ID)
}

func (l *lab) experiment(ID string) *experiment {
	l.Lock()
	defer l.Unlock()

	if e, ok := l.experiments[ID]; ok {
		return e
	}

	e := newExperiment()
	l.experiments[ID] = e
	return e
}

func (l *lab) ExperimentAim(ID string) AudienceAim {
	return l.experiment(ID).getAim()
}

func (l *lab) Session(v Visitor) Session {
	return &session{l, v}
}
