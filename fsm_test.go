package gosm_test

import (
	"context"
	"github.com/threeq/goassert"
	"github.com/threeq/gosm"
	"testing"
)

type TestEntity struct {
	id string
	s  gosm.State
}

func (t *TestEntity) ID() string {
	return t.id
}

func (t *TestEntity) State() gosm.State {
	return t.s
}

func NewTestEntity(id, state string) *TestEntity {
	return &TestEntity{id, state}
}

func testMachine() *gosm.StateMachine {

	var e1Action = func(ctx context.Context, entity gosm.Entity, from, to gosm.State) error {
		e := entity.(*TestEntity)
		e.s = to
		return nil
	}

	sm := gosm.NewStateMachine()
	sm.Name = "xxx"
	sm.Transition([]*gosm.Transition{
		{From: "s1", Event: "e1", To: "s2", Cond: gosm.Any, Action: e1Action, CondDesc: "Any"},
		{From: "s2", Event: "e2", To: "s3", Cond: gosm.Any, Action: e1Action, CondDesc: "Any"},
	}...)
	return sm
}

func TestStateMachine_Show(t *testing.T) {

	raw := testMachine().Show()
	println(raw)
	goassert.New(t).That(raw).Contains("F(s1 = e1<Any>) -> s2；")
	goassert.New(t).That(raw).Contains("F(s2 = e2<Any>) -> s3；")
}

func TestStateMachine_Trigger(t *testing.T) {
	sm := testMachine()
	entity := NewTestEntity("11","s1")
	err := sm.Trigger(context.Background(), entity, "e1")

	goassert.That(t, err).Equal(nil)
	goassert.That(t, entity.s).Equal("s2")
}
