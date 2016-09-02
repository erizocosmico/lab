package lab

import (
	"hash/crc32"
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

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// AimRandom will show the experiment randomly.
func AimRandom() AudienceAim {
	return &random{rnd}
}

type random struct {
	rnd randGenerator
}

func (r *random) shows(_ Visitor, _ StrategyGetter) bool {
	return r.rnd.Intn(101)%2 == 0
}

type percent struct {
	percent int
}

func (p *percent) shows(v Visitor, _ StrategyGetter) bool {
	return crc32.ChecksumIEEE([]byte(v.ID()))%100 < uint32(p.percent)
}

// AimPercent will randomly show the experiment to n % of the visitors.
func AimPercent(n int) AudienceAim {
	return &percent{n}
}

type strategyAim struct {
	id string
	fn func(Visitor) Params
}

func (a *strategyAim) shows(v Visitor, s StrategyGetter) bool {
	if strategy, ok := s.Strategy(a.id); ok {
		var params = make(Params)
		if a.fn != nil {
			params = a.fn(v)
		}

		return strategy(v, params)
	}
	return false
}

// Params will hold all the parameters to call a strategy.
type Params map[string]interface{}

// AimStrategy will show the experiment to the visitor if the given strategy determines
// it is ok to show it. A callback function can be given to return the Params to call the
// strategy based on the visitor.
func AimStrategy(id string, fn func(Visitor) Params) AudienceAim {
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
