package gosm

import (
	"context"
	"fmt"
	"github.com/threeq/goassert"
	"testing"
)

func TestHsm(t *testing.T) {
	m1 := NewMachine()
	m1s1 := m1.State(11)
	m1s2 := m1.State(12)
	m1e1 := Event(11)
	m1e2 := Event(12)

	m2 := NewMachine()
	s1 := State(1)
	s2 := State(2)
	s3 := State(3)
	e1 := Event(1)
	e2 := Event(2)
	e3 := Event(3)

	m2.Trans(s1.Exit(e1, Any), s2.Entry(testAction))
	m2.Trans(s2.Exit(e2, Any), s3.Entry())
	m2.Exit(s1, e3, Any).Link(Fork(Serial(All) ,s3.Entry(testAction), m1.Entry(m1s1, subMachineAction)))
	m2.Exit(s1, e2, Any).End(testAction)
	m2.Exit(s2, e1, Any).Link(m1.Entry(m1s1, Noop))

	m1.Trans(m1s1.Exit(m1e1), m1s2.Entry())
	m1.Trans(m1s1.Exit(m1e2), m1s1.Entry())
	m1.Exit(m1s2, m1e1, Any).End()
	m1.Exit(m1s2, m1e2, Any).Link(m2.Entry(s2))

	//m2.Trigger(context.Background(), NewTestEntity("1", s1), Event(1))
	err := m2.Trigger(context.Background(), NewTestEntity("1", s1), Event(3))
	goassert.That(t, err).Equal(nil)
	err = m2.Trigger(context.Background(), NewTestEntity("1", s1), e2)
	goassert.That(t, err).Equal(nil)
}

func subMachineAction(ctx context.Context, entity Entity, from IState, to IState) error {
	fmt.Printf("subMachineAction %v, %v \n", from, to)
	return nil
}

func testAction(ctx context.Context, entity Entity, from, to IState) error {
	fmt.Printf("testAction %v, %v \n", from, to)
	return nil
}
