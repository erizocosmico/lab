package lab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockVisitor struct {
	id string
}

func (m mockVisitor) ID() string {
	return m.id
}

func newVisitor(id string) Visitor {
	return mockVisitor{id}
}

func TestAimNobody(t *testing.T) {
	cases := []struct {
		visitor Visitor
		shows   bool
	}{
		{newVisitor("foo"), false},
		{newVisitor("bar"), false},
		{newVisitor("hello world"), false},
		{newVisitor("baz"), false},
	}

	aim := AimNobody()
	for _, c := range cases {
		assert.Equal(t, c.shows, aim.shows(c.visitor, nil))
	}
}

func TestAimEveryone(t *testing.T) {
	cases := []struct {
		visitor Visitor
		shows   bool
	}{
		{newVisitor("foo"), true},
		{newVisitor("bar"), true},
		{newVisitor("hello world"), true},
		{newVisitor("baz"), true},
	}

	aim := AimEveryone()
	for _, c := range cases {
		assert.Equal(t, c.shows, aim.shows(c.visitor, nil))
	}
}

func TestAimPercent(t *testing.T) {
	cases := []struct {
		visitor Visitor
		shows   bool
		percent int
	}{
		{newVisitor("hello world"), false, 50},
		{newVisitor("hello world"), true, 60},
	}

	for _, c := range cases {
		assert.Equal(t, c.shows, AimPercent(c.percent).shows(c.visitor, nil))
	}
}

type mockRandGenerator struct {
	n   []int
	pos int
}

func (m *mockRandGenerator) Intn(n int) int {
	if m.pos >= len(m.n) {
		m.pos = 0
	}
	n = m.n[m.pos]
	m.pos++
	return n
}

func newRandGenerator() randGenerator {
	return &mockRandGenerator{
		[]int{1, 3, 4, 5, 6},
		0,
	}
}

func TestAimRandom(t *testing.T) {
	aim := random{newRandGenerator()}
	cases := []bool{false, false, true, false, true}
	visitor := newVisitor("foo")
	for _, c := range cases {
		assert.Equal(t, c, aim.shows(visitor, nil))
	}
}

type strategyGetterMock struct {
}

func (strategyGetterMock) Strategy(_ string) (Strategy, bool) {
	return func(p Params) bool {
		return p["show"].(bool)
	}, true
}

func TestAimStrategy(t *testing.T) {
	visitor := newVisitor("foo")
	aim := AimStrategy("foo", func(_ Visitor) Params {
		return Params{"show": true}
	})
	assert.True(t, aim.shows(visitor, &strategyGetterMock{}))

	aim = AimStrategy("foo", func(_ Visitor) Params {
		return Params{"show": false}
	})
	assert.False(t, aim.shows(visitor, &strategyGetterMock{}))
}

func TestOr(t *testing.T) {
	aim := Or(
		AimPercent(60),
		AimStrategy("foo", func(_ Visitor) Params {
			return Params{"show": false}
		}),
	)
	assert.True(t, aim.shows(newVisitor("hello world"), &strategyGetterMock{}))

	aim = Or(
		AimPercent(50),
		AimStrategy("foo", func(_ Visitor) Params {
			return Params{"show": false}
		}),
	)
	assert.False(t, aim.shows(newVisitor("hello world"), &strategyGetterMock{}))
}

func TestAnd(t *testing.T) {
	aim := And(
		AimPercent(60),
		AimStrategy("foo", func(_ Visitor) Params {
			return Params{"show": false}
		}),
	)
	assert.False(t, aim.shows(newVisitor("hello world"), &strategyGetterMock{}))

	aim = And(
		AimPercent(60),
		AimStrategy("foo", func(_ Visitor) Params {
			return Params{"show": true}
		}),
	)
	assert.True(t, aim.shows(newVisitor("hello world"), &strategyGetterMock{}))
}
