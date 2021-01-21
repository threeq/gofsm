package gosm

import (
)

type TestEntity struct {
	id string
	s  IState
}

func (t *TestEntity) ID() string {
	return t.id
}

func (t *TestEntity) State() IState {
	return t.s
}

func NewTestEntity(id string, s  IState) *TestEntity {
	return &TestEntity{id, s}
}
//
//func testMachine() *gosm.StateMachine {
//
//	var e1Action = func(ctx context.Context, entity gosm.Entity, from, to gosm.IState) error {
//		e := entity.(*TestEntity)
//		e.s = to
//		return nil
//	}
//
//	sm := gosm.NewStateMachine()
//	sm.Name = "xxx"
//	sm.Trans([]*gosm.Trans{
//		{From: "s1", Event: "e1", To: "s2", Condition: gosm.Any, Action: e1Action, CondDesc: "Any"},
//		{From: "s2", Event: "e2", To: "s3", Condition: gosm.Any, Action: e1Action, CondDesc: "Any"},
//	}...)
//	return sm
//}
//
//func TestStateMachine_Show(t *testing.T) {
//
//	raw := testMachine().Show()
//	println(raw)
//	goassert.New(t).That(raw).Contains("F(s1 = e1<Any>) -> s2")
//	goassert.New(t).That(raw).Contains("F(s2 = e2<Any>) -> s3")
//}
//
//func TestStateMachine_Trigger(t *testing.T) {
//	sm := testMachine()
//	entity := NewTestEntity("11", "s1")
//	err := sm.Trigger(context.Background(), entity, "e1")
//
//	goassert.That(t, err).Equal(nil)
//	goassert.That(t, entity.s).Equal("s2")
//}
