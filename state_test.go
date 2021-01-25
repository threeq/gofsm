package gosm

import (
	"github.com/threeq/goassert"
	"testing"
)

func TestState_ID(t *testing.T) {
	m := NewMachine("TestState_ID")
	s1 := m.State("1")
	s2 := m.State("1")
	goassert.That(t, s1.ID()).Equal(s2.ID())

	s3 := m.State(1)
	s4 := State(1)
	goassert.That(t, s3.ID()).Equal(s4.ID())

	s5 := State("2")
	s6 := State(2)
	goassert.That(t, s5.ID()).NotEqual(s6.ID())
	goassert.That(t, s1.ID()).NotEqual(s5.ID())
}
