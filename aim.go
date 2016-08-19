package lab

import (
	"math/rand"
	"time"
)

// AudienceAim is something that determines the audience that should be shown the experiment.
type AudienceAim interface {
	// shows says if the Visitor will be in the audience that will be shown the experiment.
	shows(Visitor, StrategyGetter) bool
}

// StrategyGetter will retrieve strategies based on their ID.
type StrategyGetter interface {
	// Strategy returns the Strategy for the given ID.
	Strategy(string) (Strategy, bool)
}

type strategyGetter interface {
	getStrategies() map[string]Strategy
}

type allVisitors struct {
	show bool
}

func (a *allVisitors) shows(_ Visitor, _ StrategyGetter) bool { return a.show }

// AimNobody will show the experiment to nobody.
func AimNobody() AudienceAim {
	return &allVisitors{false}
}

// AimEveryone will show the experiment to everyone.
func AimEveryone() AudienceAim {
	return &allVisitors{true}
}

type randGenerator interface {
	Intn(int) int
}

var rnd randGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

type percent struct {
	percent int
}

func (p *percent) shows(_ Visitor, _ StrategyGetter) bool {
	return rnd.Intn(101) > p.percent
}

// AimPercent will randomly show the experiment to n % of the visitors.
func AimPercent(n int) AudienceAim {
	return &percent{n}
}

type strategyAim struct {
	id         string
	fillParams func(Params)
}

func (a *strategyAim) shows(_ Visitor, s StrategyGetter) bool {
	params := newParams()
	a.fillParams(params)
	if strategy, ok := s.Strategy(a.id); ok {
		return strategy(params)
	}
	return false
}

// AimStrategy will show the experiment to the visitor if the given strategy determines
// it is ok to show it.
func AimStrategy(id string, fn func(Params)) AudienceAim {
	return &strategyAim{id, fn}
}

type or struct {
	aims []AudienceAim
}

func (o *or) shows(v Visitor, s StrategyGetter) bool {
	for _, a := range o.aims {
		if a.shows(v, s) {
			return true
		}
	}
	return false
}

// Or shows the experiment if one or more of the aims will show the experiment.
func Or(aims ...AudienceAim) AudienceAim {
	return &or{aims}
}

type and struct {
	aims []AudienceAim
}

func (o *and) shows(v Visitor, s StrategyGetter) bool {
	for _, a := range o.aims {
		if !a.shows(v, s) {
			return false
		}
	}
	return true
}

// And will show the experiment if all of the aims will show the experiment.
func And(aims ...AudienceAim) AudienceAim {
	return &and{aims}
}
