package gosm

import (
	"context"
	"fmt"
	"github.com/threeq/goassert"
	"testing"
)

func TestStateMachine_Show(t *testing.T) {

	m1 := NewMachine()
	m1.Name = "m1"
	m1s1 := m1.State(11)
	m1s2 := m1.State(12)
	m1e1 := Event(11)
	m1e2 := Event(12)

	m2 := NewMachine()
	m2.Name = "m2"
	s1 := State(1)
	s2 := State(2)
	s3 := State(3)
	e1 := Event(1)
	e2 := Event(2)
	e3 := Event(3)

	m2.Trans(s1.Exit(e1, "A", Any), s2.Entry("", testAction))
	m2.Trans(s1.Exit(e1, "B", Any), s2.Entry("", testAction))
	m2.Trans(s1.Exit(e1, "C", Any), s2.Entry("", testAction))
	m2.Trans(s2.Exit(e2, "Any", Any), s3.Entry(""))
	m2.Fork(m2.Exit(s1, e3, "Any", Any)).Link(Serial(All), s3.Entry("", Noop), m1.Entry(m1s1, "", subMachineAction))
	m2.Exit(s1, e2, "Any", Any).End(testAction)
	m2.Exit(s3, e1, "Any", Any).Link(m1.Entry(m1s1, "", Noop))

	m1.Trans(m1s1.Exit(m1e1, "Any"), m1s2.Entry("", Noop))
	m1.Trans(m1s1.Exit(m1e2, "Any"), m1s1.Entry(""))
	m1.Exit(m1s2, m1e1, "Any", Any).End()
	m1.Exit(m1s2, m1e2, "Any", Any).Link(m2.Entry(s2, ""))

	m2.Show()

	err := m2.Trigger(context.Background(), NewTestEntity("1", s1), Event(3))
	goassert.That(t, err).Equal(nil)
	err = m2.Trigger(context.Background(), NewTestEntity("1", s1), e2)
	goassert.That(t, err).Equal(nil)

	err = m2.Trigger(context.Background(), NewTestEntity("1", State("not found")), e2)
	goassert.That(t, err.Error()).Contains("状态没有定义")

	err = m2.Trigger(context.Background(), NewTestEntity("1", s1), Event(4))
	goassert.That(t, err.Error()).Contains("事件没有定义")
}

func subMachineAction(ctx context.Context, entity Entity, from IState, to IState) error {
	fmt.Printf("subMachineAction %v, %v \n", from, to)
	return nil
}

func testAction(ctx context.Context, entity Entity, from, to IState) error {
	fmt.Printf("testAction %v, %v \n", from, to)
	return nil
}

func testMachine() *StateMachine {

	var e1Action = func(ctx context.Context, entity Entity, from, to IState) error {
		e := entity.(*TestEntity)
		e.s = to
		return nil
	}

	sm := NewMachine()
	sm.Name = "xxx"
	sm.Trans(State("s1").Exit("e1", ""), State("s2").Entry("", e1Action))
	sm.Trans(State("s2").Exit("e2", "Any", Any), State("s3").Entry("", e1Action))
	sm.Trans(State("s3").Exit("e2", "always", Any), State("s3").Entry("", e1Action))
	sm.Trans(State("s3").Exit("e3", "", Any), State("s3").Entry("", e1Action))
	return sm
}

func TestStateMachine_Trigger(t *testing.T) {
	sm := testMachine()
	entity := NewTestEntity("11", State("s1"))
	err := sm.Trigger(context.Background(), entity, "e1")

	goassert.That(t, err).Equal(nil)
	goassert.That(t, entity.s.ID()).Equal("s2")
}
