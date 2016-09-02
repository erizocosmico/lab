package lab

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLab(t *testing.T) {
	lab := New()
	lab.DefineStrategy("is-admin", func(p Params) bool {
		v, ok := p["admin"].(bool)
		if !ok {
			return false
		}
		return v
	})

	lab.Experiment("secret-stuff").
		Aim(AimPercent(70)).
		Aim(AimStrategy("is-admin", func(v Visitor) Params {
			return Params{
				"admin": strings.HasPrefix(v.ID(), "admin"),
			}
		}))

	cases := []struct {
		visitor  Visitor
		launched bool
	}{
		{newVisitor("admin"), true},
		{newVisitor("admin-bazzinga"), false},
		{newVisitor("bazzinga"), false},
		{newVisitor("hello world"), false},
	}

	var called int
	for _, c := range cases {
		session := lab.Session(c.visitor)
		ok := session.Launch("secret-stuff", func() {
			called++
		})
		assert.Equal(t, c.launched, ok)
		assert.False(t, session.Launch("nope", nil))
	}
	assert.Equal(t, 1, called)
}
