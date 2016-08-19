package lab

import "sync"

// Experiment is a single experiment that wants to be run.
type Experiment interface {
	// Aim defines a restriction over the audience that will be shown this experiment.
	Aim(AudienceAim) Experiment
}

type experiment struct {
	sync.RWMutex
	aim AudienceAim
}

func newExperiment() *experiment {
	return &experiment{aim: nil}
}

func (e *experiment) Aim(aim AudienceAim) Experiment {
	e.Lock()
	defer e.Unlock()

	if e.aim == nil {
		e.aim = aim
	} else {
		And(e.aim, aim)
	}
	return e
}

func (e *experiment) getAim() AudienceAim {
	e.Lock()
	defer e.Unlock()
	return e.aim
}
