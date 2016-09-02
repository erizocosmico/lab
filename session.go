package lab

// Session is a single session for a visitor.
type Session interface {
	// Launch will run the given function for the given experiment ID only if the
	// visitor is a target of the experiment. Will return as well a boolean to report
	// if the visitor was shown the experiment.
	Launch(string, func()) bool
}

type session struct {
	lab     *lab
	visitor Visitor
}

func (s *session) Launch(ID string, fn func()) bool {
	aim := s.lab.ExperimentAim(ID)
	if aim == nil {
		return false
	}

	if aim.shows(s.visitor, s.lab) {
		if fn != nil {
			fn()
		}
		return true
	}
	return false
}
